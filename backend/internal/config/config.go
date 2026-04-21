package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// ConnectDBFromEnv reads DB_DSN from environment
// DSN format: user:pass@tcp(host:3306)/dbname?parseTime=true
func ConnectDBFromEnv() (*sql.DB, error) {
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        return nil, fmt.Errorf("DB_DSN not set")
    }
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    maxOpen := 10
    if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
        fmt.Sscanf(v, "%d", &maxOpen)
    }
    db.SetMaxOpenConns(maxOpen)
    db.SetConnMaxIdleTime(5 * time.Minute)
    if err := db.Ping(); err != nil {
        return nil, err
    }
    log.Printf("connected to db, maxOpen=%d", maxOpen)
    return db, nil
}
