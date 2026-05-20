// internal/config/env.go

package config

import (
	"fmt"
	"os"
)

// =========================
// env accessors
// =========================

func AppEnv() string {
	return os.Getenv("APP_ENV")
}

func RunMode() string {
	return os.Getenv("RUN_MODE")
}

func DevAuthBypassEnabled() bool {

	return os.Getenv(
		"ENABLE_DEV_AUTH_BYPASS",
	) == "true"
}

// =========================
// runtime helpers
// =========================

func IsProduction() bool {
	return AppEnv() == "production"
}

func IsLocal() bool {
	return RunMode() == "local"
}

func IsLambda() bool {
	return RunMode() == "lambda"
}

// =========================
// startup validation
// =========================

func ValidateEnv() error {

	// =========================
	// production guard
	// =========================

	if IsProduction() {

		if DevAuthBypassEnabled() {

			return fmt.Errorf(
				"ENABLE_DEV_AUTH_BYPASS=true is forbidden in production",
			)
		}

		if IsLocal() {

			return fmt.Errorf(
				"RUN_MODE=local is forbidden in production",
			)
		}
	}

	// =========================
	// dev auth guard
	// =========================

	if DevAuthBypassEnabled() &&
		!IsLocal() {

		return fmt.Errorf(
			"ENABLE_DEV_AUTH_BYPASS requires RUN_MODE=local",
		)
	}

	return nil
}
