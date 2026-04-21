package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"portfolio/backend/internal/config"
	"portfolio/backend/internal/handlers"
	"portfolio/backend/internal/repository"
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
	userHandler := handlers.NewUserHandler(userSvc)

	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/users", userHandler.HandleUsers)
	mux.HandleFunc("/api/v1/users/", userHandler.HandleUserByID)
	mux.HandleFunc("/api/v1/tasks", taskHandler.HandleTasks)
	mux.HandleFunc("/api/v1/tasks/", taskHandler.HandleTaskByID)

	// Serve swagger YAML
	mux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))

	// Serve docs via template
	mux.HandleFunc("/api/v1/docs", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./swagger/template/index.html.tmpl"))
		data := map[string]string{
			"SchemaURL": "/api/v1/spec/swagger.yml",
			"DomID":     "#root",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	})

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
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not initialized", http.StatusInternalServerError)
	})
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not initialized", http.StatusInternalServerError)
	})
	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not initialized", http.StatusInternalServerError)
	})
	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not initialized", http.StatusInternalServerError)
	})

	// Serve swagger YAML
	mux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))

	// Serve docs via template
	mux.HandleFunc("/api/v1/docs", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./swagger/template/index.html.tmpl"))
		data := map[string]string{
			"SchemaURL": "/api/v1/spec/swagger.yml",
			"DomID":     "#root",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	})

	initHandler(mux)
	return mux
}

func initHandler(mux *http.ServeMux) {
	db, err := config.ConnectDBFromEnv()
	if err != nil {
		log.Printf("warning: failed to connect DB in initHandler: %v", err)
		return
	}
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	mux.HandleFunc("/api/v1/users", userHandler.HandleUsers)
	mux.HandleFunc("/api/v1/users/", userHandler.HandleUserByID)
	mux.HandleFunc("/api/v1/tasks", taskHandler.HandleTasks)
	mux.HandleFunc("/api/v1/tasks/", taskHandler.HandleTaskByID)
}

func proxyToHTTP(ctx context.Context, handler http.Handler, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := req.RequestContext.HTTP.Method
	rawPath := req.RawPath
	if rawPath == "" {
		rawPath = req.RequestContext.HTTP.Path
	}
	u := url.URL{
		Path:     rawPath,
		RawQuery: req.RawQueryString,
	}
	var body io.ReadCloser
	if req.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{StatusCode: 400}, nil
		}
		body = io.NopCloser(bytes.NewReader(decoded))
	} else {
		body = io.NopCloser(strings.NewReader(req.Body))
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, nil
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	if host, ok := req.Headers["host"]; ok {
		httpReq.Host = host
	}

	rw := &responseRecorder{
		HeaderMap: http.Header{},
		Body:      &bytes.Buffer{},
	}
	handler.ServeHTTP(rw, httpReq)

	status := rw.Status
	if status == 0 {
		status = http.StatusOK
	}

	resp := events.APIGatewayV2HTTPResponse{
		StatusCode:      status,
		Headers:         map[string]string{},
		Body:            rw.Body.String(),
		IsBase64Encoded: false,
	}
	for k, vv := range rw.HeaderMap {
		if len(vv) > 0 {
			resp.Headers[k] = vv[0]
		}
	}
	return resp, nil
}

type responseRecorder struct {
	HeaderMap http.Header
	Body      *bytes.Buffer
	Status    int
}

func (r *responseRecorder) Header() http.Header { return r.HeaderMap }
func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.Status == 0 {
		r.Status = http.StatusOK
	}
	return r.Body.Write(b)
}
func (r *responseRecorder) WriteHeader(status int) { r.Status = status }

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
