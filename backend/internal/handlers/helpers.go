package handlers

import (
	"net/http"
	"time"

	"portfolio/backend/internal/auth"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/validator"
)

// =========================
// parseID
// =========================

func parseID(
	w http.ResponseWriter,
	r *http.Request,
) (int64, bool) {

	id, err := httpx.PathID(r)

	if err != nil {
		httpx.WriteError(
			w,
			http.StatusBadRequest,
			httpx.CodeInvalidID,
			"invalid id",
		)
		return 0, false
	}

	return id, true
}

// =========================
// decodeAndValidate
// =========================

func decodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) bool {

	if err := httpx.DecodeJSON(w, r, dst); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			httpx.CodeInvalidJSON,
			err.Error(),
		)

		return false
	}

	if errs := validator.ValidateStruct(dst); errs != nil {

		httpx.WriteValidationErrors(w, errs)
		return false
	}

	return true
}

// =========================
// requireAuthUserID (FIXED)
// =========================

func requireAuthUserID(
	w http.ResponseWriter,
	r *http.Request,
) (int64, bool) {

	userID, ok := auth.GetUserID(r.Context())

	if !ok || userID <= 0 {

		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			httpx.CodeUnauthorized,
			"missing user context",
		)

		return 0, false
	}

	return userID, true
}

// =========================
// parseOptionalDueDate
// =========================

func parseOptionalDueDate(
	w http.ResponseWriter,
	value string,
) (*time.Time, bool) {

	t, err := httpx.ParseOptionalTime(value)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			httpx.CodeInvalidDueDate,
			"invalid due_date format",
		)

		return nil, false
	}

	return t, true
}
