package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"portfolio/backend/internal/apierrors"
)

// =========================
// JSON response
// =========================

func WriteJSON(
	w http.ResponseWriter,
	status int,
	v any,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {

		log.Printf(
			"[ERROR] response encode failed: %v",
			err,
		)
	}
}

// =========================
// error response
// =========================

func WriteError(
	w http.ResponseWriter,
	status int,
	code string,
	msg string,
) {

	WriteJSON(
		w,
		status,
		apierrors.ErrorResponse{
			Code:    code,
			Message: msg,
		},
	)
}

// =========================
// validation error response
// =========================

func WriteValidationErrors(
	w http.ResponseWriter,
	errs map[string]string,
) {

	WriteJSON(
		w,
		http.StatusBadRequest,
		apierrors.ValidationErrorResponse{
			Code:    apierrors.CodeValidationError,
			Message: "validation failed",
			Errors:  errs,
		},
	)
}
