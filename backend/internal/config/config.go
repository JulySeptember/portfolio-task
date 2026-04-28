package config

import (
        "database/sql"
        "fmt"
        "log"
        "os"
        "time"

        _ "github.com/go-sql-driver/mysql"
)

// ConnectDBFromEnv reads DSN and pool settings from environment and returns *sql.DB.
// Environment variables:
//   - DB_DSN (required)
//   - DB_MAX_OPEN_CONNS (optional, default 10)
//   - DB_MAX_IDLE_CONNS (optional, default = DB_MAX_OPEN_CONNS)
//   - DB_CONN_MAX_LIFETIME (optional, duration string like "5m", default 5m)
func ConnectDBFromEnv() (*sql.DB, error) {
        dsn := os.Getenv("DB_DSN")
        if dsn == "" {
                return nil, fmt.Errorf("DB_DSN not set")
        }
        db, err := sql.Open("mysql", dsn)
        if err != nil {
                return nil, err
        }

        // Max open conns
        maxOpen := 10
        if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
                fmt.Sscanf(v, "%d", &maxOpen)
        }
        db.SetMaxOpenConns(maxOpen)

        // Max idle conns
        maxIdle := maxOpen
        if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
                fmt.Sscanf(v, "%d", &maxIdle)
        }
        db.SetMaxIdleConns(maxIdle)

        // Conn max lifetime
        connMaxLifetime := 5 * time.Minute
        if v := os.Getenv("DB_CONN_MAX_LIFETIME"); v != "" {
                if d, err := time.ParseDuration(v); err == nil {
                        connMaxLifetime = d
                }
        }
        db.SetConnMaxLifetime(connMaxLifetime)

        if err := db.Ping(); err != nil {
                return nil, err
        }
        log.Printf("connected to db, maxOpen=%d, maxIdle=%d, connMaxLifetime=%v", maxOpen, maxIdle, connMaxLifetime)
        return db, nil
}
