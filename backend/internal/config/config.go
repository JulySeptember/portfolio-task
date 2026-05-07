package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDBFromEnv reads DSN and pool settings from environment and returns *sql.DB.
// Environment variables:
//   - DB_DSN (required)
//   - DB_MAX_OPEN_CONNS (optional, default 10)
//   - DB_MAX_IDLE_CONNS (optional, default = DB_MAX_OPEN_CONNS)
//   - DB_CONN_MAX_LIFETIME (optional, duration string like "5m", default 5m)
//   - DB_PING_TIMEOUT (optional, duration string like "3s", default 3s)
func ConnectDBFromEnv() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// parse max open conns
	maxOpen := 10
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxOpen = n
		} else {
			log.Printf("invalid DB_MAX_OPEN_CONNS=%q, using default %d", v, maxOpen)
		}
	}
	db.SetMaxOpenConns(maxOpen)

	// parse max idle conns (default = maxOpen)
	maxIdle := maxOpen
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			maxIdle = n
		} else {
			log.Printf("invalid DB_MAX_IDLE_CONNS=%q, using default %d", v, maxIdle)
		}
	}
	db.SetMaxIdleConns(maxIdle)

	// parse conn max lifetime (duration string)
	connMaxLifetime := 5 * time.Minute
	if v := os.Getenv("DB_CONN_MAX_LIFETIME"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			connMaxLifetime = d
		} else {
			log.Printf("invalid DB_CONN_MAX_LIFETIME=%q, using default %v", v, connMaxLifetime)
		}
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	// ping timeout (duration string)
	pingTimeout := 3 * time.Second
	if v := os.Getenv("DB_PING_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			pingTimeout = d
		} else {
			log.Printf("invalid DB_PING_TIMEOUT=%q, using default %v", v, pingTimeout)
		}
	}

	// Ping with timeout and close db on failure to avoid leaks
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	log.Printf("connected to db, maxOpen=%d, maxIdle=%d, connMaxLifetime=%v", maxOpen, maxIdle, connMaxLifetime)
	return db, nil
}
