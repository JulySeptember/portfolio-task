// internal/config/config.go

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

// =========================
// env helpers
// =========================

func getEnvInt(
	key string,
	def int,
) int {

	v := os.Getenv(key)

	if v == "" {
		return def
	}

	n, err := strconv.Atoi(v)

	if err != nil || n < 0 {

		log.Printf(
			"invalid %s=%q, using default=%d",
			key,
			v,
			def,
		)

		return def
	}

	return n
}

func getEnvDuration(
	key string,
	def time.Duration,
) time.Duration {

	v := os.Getenv(key)

	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)

	if err != nil {

		log.Printf(
			"invalid %s=%q, using default=%v",
			key,
			v,
			def,
		)

		return def
	}

	return d
}

// =========================
// safety validation
// =========================

func validateEnvironment() {

	appEnv := os.Getenv(
		"APP_ENV",
	)

	enableDevAuthBypass :=
		os.Getenv(
			"ENABLE_DEV_AUTH_BYPASS",
		) == "true"

	// =========================
	// prevent prod auth bypass
	// =========================

	if enableDevAuthBypass &&
		appEnv != "development" {

		panic(
			"ENABLE_DEV_AUTH_BYPASS is only allowed in development",
		)
	}
}

// =========================
// ConnectDBFromEnv
// =========================

func ConnectDBFromEnv() (*sql.DB, error) {

	// =========================
	// environment validation
	// =========================

	validateEnvironment()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbHost == "" {
		return nil, fmt.Errorf("DB_HOST not set")
	}

	if dbPort == "" {
		dbPort = "3306"
	}

	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME not set")
	}

	if dbUser == "" {
		return nil, fmt.Errorf("DB_USER not set")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	db, err := sql.Open(
		"mysql",
		dsn,
	)

	if err != nil {
		return nil, err
	}

	// =========================
	// runtime mode
	// =========================

	mode := "local"

	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""

	if isLambda {
		mode = "lambda"
	}
	// =========================
	// defaults
	// =========================

	var (
		maxOpen         int
		maxIdle         int
		connMaxLifetime time.Duration
		connMaxIdleTime time.Duration
		pingTimeout     time.Duration
	)

	// =========================
	// lambda optimized
	// =========================

	if isLambda {

		// Lambda:
		// keep pool VERY small

		maxOpen = 1
		maxIdle = 1

		// recycle frozen connections

		connMaxLifetime =
			2 * time.Minute

		connMaxIdleTime =
			1 * time.Minute

		pingTimeout =
			3 * time.Second

	} else {

		// local/dev defaults

		maxOpen = 10
		maxIdle = 10

		connMaxLifetime =
			5 * time.Minute

		connMaxIdleTime =
			3 * time.Minute

		pingTimeout =
			3 * time.Second
	}

	// =========================
	// env override
	// =========================

	maxOpen = getEnvInt(
		"DB_MAX_OPEN_CONNS",
		maxOpen,
	)

	maxIdle = getEnvInt(
		"DB_MAX_IDLE_CONNS",
		maxIdle,
	)

	connMaxLifetime =
		getEnvDuration(
			"DB_CONN_MAX_LIFETIME",
			connMaxLifetime,
		)

	connMaxIdleTime =
		getEnvDuration(
			"DB_CONN_MAX_IDLE_TIME",
			connMaxIdleTime,
		)

	pingTimeout =
		getEnvDuration(
			"DB_PING_TIMEOUT",
			pingTimeout,
		)

	// =========================
	// safety guards
	// =========================

	if maxOpen <= 0 {
		maxOpen = 1
	}

	if maxIdle < 0 {
		maxIdle = 0
	}

	if maxIdle > maxOpen {
		maxIdle = maxOpen
	}

	// =========================
	// pool settings
	// =========================

	db.SetMaxOpenConns(
		maxOpen,
	)

	db.SetMaxIdleConns(
		maxIdle,
	)

	db.SetConnMaxLifetime(
		connMaxLifetime,
	)

	db.SetConnMaxIdleTime(
		connMaxIdleTime,
	)

	// =========================
	// ping
	// =========================

	ctx, cancel := context.WithTimeout(
		context.Background(),
		pingTimeout,
	)

	defer cancel()

	if err := db.PingContext(
		ctx,
	); err != nil {

		_ = db.Close()

		return nil, fmt.Errorf(
			"db ping failed: %w",
			err,
		)
	}

	log.Printf(
		"connected to db "+
			"(mode=%s maxOpen=%d maxIdle=%d lifetime=%v idleTime=%v)",
		mode,
		maxOpen,
		maxIdle,
		connMaxLifetime,
		connMaxIdleTime,
	)

	return db, nil
}
