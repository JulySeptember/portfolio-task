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
	"portfolio/backend/internal/middleware"
	"portfolio/backend/internal/repository"
	"portfolio/backend/internal/router"
	"portfolio/backend/internal/service"
	"portfolio/backend/swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// ----------------------------
// globals (lambda warm start)
// ----------------------------

var (
	globalDB      *sql.DB
	globalApp     *App
	globalHandler *httpadapter.HandlerAdapter
)

// ----------------------------
// App
// ----------------------------

type App struct {
	DB   *sql.DB
	HTTP http.Handler
}

// ----------------------------
// NewApp
// ----------------------------

func NewApp(db *sql.DB) *App {

	// =========================
	// Repository
	// =========================

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// =========================
	// Service
	// =========================

	userSvc := service.NewUserService(userRepo)
	taskSvc := service.NewTaskService(taskRepo)

	// =========================
	// Handler
	// =========================

	userHandler := handlers.NewUserHandler(userSvc)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	// =========================
	// API Router
	// =========================

	apiRouter := router.NewRouter(
		userHandler,
		taskHandler,
	)

	// JWT only for API
	apiHandler := middleware.Chain(
		apiRouter,
		middleware.JWT,
	)

	// =========================
	// Swagger
	// =========================

	swaggerHandler := swagger.Handler()

	// =========================
	// Root mux
	// =========================

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", apiHandler)

	mux.Handle(
		"/api/v1/docs/",
		swaggerHandler,
	)

	mux.Handle(
		"/api/v1/spec/",
		swaggerHandler,
	)

	// =========================
	// Global middleware
	// =========================

	handler := middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.Logging,
		middleware.CORS,
	)

	return &App{
		DB:   db,
		HTTP: handler,
	}
}

// ----------------------------
// lambda handler
// ----------------------------

func lambdaHandler() *httpadapter.HandlerAdapter {

	if globalHandler != nil {
		return globalHandler
	}

	// DB init
	if globalDB == nil {

		db, err := config.ConnectDBFromEnv()
		if err != nil {
			log.Fatal(err)
		}

		globalDB = db
	}

	// App init
	globalApp = NewApp(globalDB)

	// Lambda adapter
	globalHandler = httpadapter.New(
		globalApp.HTTP,
	)

	return globalHandler
}

// ----------------------------
// main
// ----------------------------

func main() {

	// =========================
	// local mode
	// =========================

	if os.Getenv("RUN_MODE") == "local" {

		db, err := config.ConnectDBFromEnv()
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

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

// ----------------------------
// local server
// ----------------------------

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
	}

	// graceful shutdown

	idleConnsClosed := make(chan struct{})

	go func() {

		c := make(chan os.Signal, 1)

		signal.Notify(
			c,
			syscall.SIGINT,
			syscall.SIGTERM,
		)

		<-c

		log.Println("shutting down server...")

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

		if globalDB != nil {
			globalDB.Close()
		}

		close(idleConnsClosed)
	}()

	log.Printf(
		"server started :%s",
		port,
	)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {

		log.Fatalf(
			"server failed: %v",
			err,
		)
	}

	<-idleConnsClosed
}
