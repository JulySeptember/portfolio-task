package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portfolio/backend/internal/config"
	"portfolio/backend/internal/container"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/middleware"
	"portfolio/backend/internal/router"
	"portfolio/backend/swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type App struct {
	Container *container.Container
	HTTP      http.Handler
}

func NewApp(
	c *container.Container,
) *App {

	// =========================
	// api router
	// =========================

	apiRouter := router.NewRouter(
		c.UserHandler,
		c.TaskHandler,
	)

	// =========================
	// protected api middleware
	// =========================

	apiHandler := middleware.Chain(
		apiRouter,
		middleware.AuthMiddleware,
		middleware.Logging,
	)

	// =========================
	// root mux
	// =========================

	mux := http.NewServeMux()

	// =========================
	// health
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
		middleware.CORS,
		middleware.Recovery,
	)

	return &App{
		Container: c,
		HTTP:      handler,
	}
}

// =========================
// lambda cache
// =========================

var cachedLambdaAdapter *httpadapter.HandlerAdapterV2

func lambdaAdapter() *httpadapter.HandlerAdapterV2 {

	if cachedLambdaAdapter != nil {
		return cachedLambdaAdapter
	}

	db, err := config.ConnectDBFromEnv()

	if err != nil {
		log.Fatal(err)
	}

	c, err := container.New(
		db,
	)

	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(c)

	cachedLambdaAdapter = httpadapter.NewV2(
		app.HTTP,
	)

	return cachedLambdaAdapter
}

func main() {

	if err := config.ValidateEnv(); err != nil {
		log.Fatal(err)
	}

	// =========================
	// local mode
	// =========================

	if os.Getenv("RUN_MODE") == "local" {

		db, err := config.ConnectDBFromEnv()

		if err != nil {
			log.Fatal(err)
		}

		c, err := container.New(
			db,
		)

		if err != nil {
			log.Fatal(err)
		}

		app := NewApp(c)

		runLocal(app)

		return
	}

	// =========================
	// lambda mode
	// =========================

	lambda.Start(
		func(
			ctx context.Context,
			req events.APIGatewayV2HTTPRequest,
		) (
			events.APIGatewayV2HTTPResponse,
			error,
		) {

			// =========================
			// jwt claims
			// =========================

			claims := req.RequestContext.
				Authorizer.
				JWT.
				Claims

			// =========================
			// inject headers
			// =========================

			if req.Headers == nil {
				req.Headers = map[string]string{}
			}

			if sub, ok := claims["sub"]; ok {
				req.Headers["X-Auth-Sub"] = sub
			}

			if email, ok := claims["email"]; ok {
				req.Headers["X-Auth-Email"] = email
			}

			// =========================
			// proxy
			// =========================

			return lambdaAdapter().
				ProxyWithContext(
					ctx,
					req,
				)
		},
	)
}

// =========================
// local server
// =========================

func runLocal(
	app *App,
) {

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

		if app.Container != nil &&
			app.Container.DB != nil {

			log.Println(
				"closing database connection...",
			)

			if err := app.Container.DB.Close(); err != nil {

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
