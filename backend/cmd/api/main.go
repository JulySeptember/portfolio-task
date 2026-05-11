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

var (
	globalDB      *sql.DB
	globalApp     *App
	globalHandler *httpadapter.HandlerAdapter
)

type App struct {
	DB   *sql.DB
	HTTP http.Handler
}

func NewApp(db *sql.DB) *App {

	// repository

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// service

	userSvc := service.NewUserService(userRepo)
	taskSvc := service.NewTaskService(taskRepo)

	// handler

	userHandler := handlers.NewUserHandler(userSvc)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	// api router

	apiRouter := router.NewRouter(
		userHandler,
		taskHandler,
	)

	// protected api

	apiHandler := middleware.Chain(
		apiRouter,
		middleware.JWT,
	)

	// root mux

	mux := http.NewServeMux()

	// health

	mux.HandleFunc(
		"/health",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		},
	)

	// swagger

	mux.HandleFunc(
		"/api/v1/docs/",
		swagger.DocsHandler,
	)

	mux.HandleFunc(
		"/api/v1/spec/swagger.yml",
		swagger.SpecHandler,
	)

	// api

	mux.Handle(
		"/api/v1/",
		apiHandler,
	)

	// middleware

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

	globalApp = NewApp(globalDB)

	globalHandler = httpadapter.New(
		globalApp.HTTP,
	)

	return globalHandler
}

func main() {

	// local

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

	// lambda

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

		if app.DB != nil {
			app.DB.Close()
		}

		close(idleConnsClosed)
	}()

	log.Printf(
		"server started on :%s",
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
