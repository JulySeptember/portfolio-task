package main

import (
    "database/sql"
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

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
    _ "github.com/go-sql-driver/mysql"
)

var (
    // globalDB is initialized once on cold start (Lambda) or in runLocal (dev).
    globalDB *sql.DB

    // mux is the http.Handler used by Lambda proxy and local server.
    mux http.Handler
)

func init() {
    if os.Getenv("RUN_MODE") == "local" {
        return
    }
    if err := initDB(); err != nil {
        log.Fatalf("failed to init db: %v", err)
    }
    mux = loggingMiddleware(buildMux())
}

func main() {
    if os.Getenv("RUN_MODE") == "local" {
        runLocal()
        return
    }
    adapter := httpadapter.New(mux)
    lambda.Start(adapter.ProxyWithContext)
}

func initDB() error {
    if globalDB != nil {
        return nil
    }
    db, err := config.ConnectDBFromEnv()
    if err != nil {
        return err
    }
    globalDB = db
    return nil
}

func CloseDB() error {
    if globalDB == nil {
        return nil
    }
    err := globalDB.Close()
    globalDB = nil
    return err
}

func runLocal() {
    db, err := config.ConnectDBFromEnv()
    if err != nil {
        log.Fatalf("failed to connect db: %v", err)
    }
    defer func() { _ = db.Close() }()

    globalDB = db

    userRepo := repository.NewUserRepository(db)
    userSvc := service.NewUserService(userRepo)

    taskRepo := repository.NewTaskRepo(db)
    taskSvc := service.NewTaskService(taskRepo)

    localMux := http.NewServeMux()
    router.RegisterUserRoutes(localMux, userSvc)
    router.RegisterTaskRoutes(localMux, taskSvc)

    localMux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))
    localMux.HandleFunc("/api/v1/docs", serveDocs)

    addr := ":" + getEnv("PORT", "8080")
    srv := &http.Server{
        Addr:         addr,
        Handler:      loggingMiddleware(localMux),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Printf("local server listening on %s (RUN_MODE=%s)", addr, os.Getenv("RUN_MODE"))
    log.Fatal(srv.ListenAndServe())
}

func buildMux() http.Handler {
    mux := http.NewServeMux()

    mux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))
    mux.HandleFunc("/api/v1/docs", serveDocs)

    userRepo := repository.NewUserRepository(globalDB)
    userSvc := service.NewUserService(userRepo)

    taskRepo := repository.NewTaskRepo(globalDB)
    taskSvc := service.NewTaskService(taskRepo)

    router.RegisterUserRoutes(mux, userSvc)
    router.RegisterTaskRoutes(mux, taskSvc)

    return mux
}

func serveDocs(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("./swagger/template/index.html.tmpl")
    if err != nil {
        http.Error(w, "docs not available", http.StatusInternalServerError)
        return
    }
    data := map[string]string{
        "SchemaURL": "/api/v1/spec/swagger.yml",
        "DomID":     "#root",
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    _ = tmpl.Execute(w, data)
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

        remote := r.RemoteAddr
        if remote == "" {
            remote = "-"
        }
        if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
            parts := strings.Split(xff, ",")
            remote = strings.TrimSpace(parts[0])
        }

        log.Printf("started %s %s from %s", r.Method, r.URL.Path, remote)

        next.ServeHTTP(w, r)

        log.Printf("completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
    })
}
