package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portfolio/backend/internal/config"
	"portfolio/backend/internal/handlers"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/middleware"
	"portfolio/backend/internal/repository"
	"portfolio/backend/internal/router"
	"portfolio/backend/internal/service"
	"portfolio/backend/swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var (
	globalDB      *sql.DB
	globalHandler *httpadapter.HandlerAdapter
)

type App struct {
	DB   *sql.DB
	HTTP http.Handler
}

func NewApp(db *sql.DB) *App {

	// =========================
	// repository
	// =========================

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// =========================
	// service
	// =========================

	userSvc := service.NewUserService(userRepo)
	taskSvc := service.NewTaskService(taskRepo)

	// =========================
	// handler
	// =========================

	userHandler := handlers.NewUserHandler(userSvc)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	// =========================
	// router
	// =========================

	apiRouter := router.NewRouter(
		userHandler,
		taskHandler,
	)

	// =========================
	// api middleware
	// =========================
	// execution order:
	// Auth -> Logging -> Router
	//
	// Logging can access user_id
	// after AuthMiddleware sets context
	// =========================

	apiHandler := middleware.Chain(
		apiRouter,
		middleware.Logging,
		middleware.AuthMiddleware(userSvc),
	)

	// =========================
	// root mux
	// =========================

	mux := http.NewServeMux()

	// =========================
	// public endpoints
	// =========================

	mux.HandleFunc(
		"/health",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			httpx.WriteJSON(
				w,
				http.StatusOK,
				map[string]string{
					"status": "ok",
				},
			)
		},
	)

	// =========================
	// swagger
	// =========================

	mux.HandleFunc(
		"/api/v1/docs/",
		swagger.DocsHandler,
	)

	mux.HandleFunc(
		"/api/v1/spec/swagger.yml",
		swagger.SpecHandler,
	)

	// =========================
	// protected api
	// =========================

	mux.Handle(
		"/api/v1/",
		apiHandler,
	)

	// =========================
	// global middleware
	// =========================

	handler := middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.CORS,
	)

	return &App{
		DB:   db,
		HTTP: handler,
	}
}

// =========================
// lambda singleton
// =========================

func lambdaHandler() *httpadapter.HandlerAdapter {

	if globalHandler != nil {
		return globalHandler
	}

	if globalDB == nil {

		db, err := config.ConnectDBFromEnv()
		if err != nil {
			log.Fatal(err)
		}

		globalDB = db
	}

	app := NewApp(globalDB)

	globalHandler = httpadapter.New(
		app.HTTP,
	)

	return globalHandler
}

func main() {

	// =========================
	// local mode
	// =========================

	if os.Getenv("RUN_MODE") == "local" {

		db, err := config.ConnectDBFromEnv()
		if err != nil {
			log.Fatal(err)
		}

		app := NewApp(db)

		runLocal(app)

		return
	}

	// =========================
	// lambda mode
	// =========================

	lambda.Start(
		func(
			ctx context.Context,
			req events.APIGatewayProxyRequest,
		) (
			events.APIGatewayProxyResponse,
			error,
		) {

			return lambdaHandler().
				ProxyWithContext(ctx, req)
		},
	)
}

// =========================
// local server
// =========================

func runLocal(app *App) {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      app.HTTP,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdownDone := make(chan struct{})

	go func() {

		sigCh := make(chan os.Signal, 1)

		signal.Notify(
			sigCh,
			syscall.SIGINT,
			syscall.SIGTERM,
		)

		<-sigCh

		log.Println(
			"shutting down server...",
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)

		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {

			log.Printf(
				"server shutdown error: %v",
				err,
			)
		}

		if app.DB != nil {

			log.Println(
				"closing database connection...",
			)

			if err := app.DB.Close(); err != nil {

				log.Printf(
					"db close error: %v",
					err,
				)
			}
		}

		close(shutdownDone)
	}()

	log.Printf(
		"server started on :%s",
		port,
	)

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {

		log.Fatalf(
			"server failed: %v",
			err,
		)
	}

	<-shutdownDone
}
