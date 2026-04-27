package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"portfolio/backend/internal/config"
	"portfolio/backend/internal/repository"
	"portfolio/backend/internal/router"
	"portfolio/backend/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if os.Getenv("RUN_MODE") == "local" {
		runLocal()
		return
	}

	mux := buildMux()

	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return proxyToHTTP(ctx, mux, req)
	})
}

func runLocal() {
	db, err := config.ConnectDBFromEnv()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo)

	mux := http.NewServeMux()

	// Router 層を使用
	router.RegisterUserRoutes(mux, userSvc)
	router.RegisterTaskRoutes(mux, taskSvc)

	// Swagger
	mux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))
	mux.HandleFunc("/api/v1/docs", serveDocs)

	addr := ":" + getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("local server listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func buildMux() http.Handler {
	mux := http.NewServeMux()

	// Swagger
	mux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))
	mux.HandleFunc("/api/v1/docs", serveDocs)

	initHandler(mux)
	return mux
}

func initHandler(mux *http.ServeMux) {
	db, err := config.ConnectDBFromEnv()
	if err != nil {
		log.Fatalf("failed to connect DB in initHandler: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo)

	// Router 層を使用
	router.RegisterUserRoutes(mux, userSvc)
	router.RegisterTaskRoutes(mux, taskSvc)
}

func serveDocs(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./swagger/template/index.html.tmpl"))
	data := map[string]string{
		"SchemaURL": "/api/v1/spec/swagger.yml",
		"DomID":     "#root",
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, data)
}

func proxyToHTTP(ctx context.Context, mux http.Handler, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	r, err := http.NewRequest(req.RequestContext.HTTP.Method, req.RawPath, strings.NewReader(req.Body))
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	// ヘッダーをコピー
	for k, v := range req.Headers {
		r.Header.Set(k, v)
	}

	// クエリパラメータをコピー
	q := r.URL.Query()
	for k, v := range req.QueryStringParameters {
		q.Set(k, v)
	}
	r.URL.RawQuery = q.Encode()

	// レスポンスをキャプチャ
	w := &responseCapture{header: http.Header{}}
	mux.ServeHTTP(w, r)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: w.status,
		Headers:    convertHeaders(w.header),
		Body:       w.body.String(),
	}, nil
}

type responseCapture struct {
	header http.Header
	body   strings.Builder
	status int
}

func (r *responseCapture) Header() http.Header {
	return r.header
}

func (r *responseCapture) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseCapture) WriteHeader(statusCode int) {
	r.status = statusCode
}

func convertHeaders(h http.Header) map[string]string {
	out := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(w, r)

		log.Printf("completed in %v", time.Since(start))
	})
}
