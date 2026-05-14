package httpx

import (
	"encoding/json"
	"net/http"
)

// =========================
// JSON
// =========================

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

// =========================
// Error
// =========================

func WriteError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {

	WriteJSON(
		w,
		status,
		ErrorResponse{
			Code:    code,
			Message: message,
		},
	)
}

// =========================
// Validation Error
// =========================

func WriteValidationErrors(
	w http.ResponseWriter,
	fields map[string]string,
) {

	WriteJSON(
		w,
		http.StatusBadRequest,
		ValidationErrorResponse{
			Code:    CodeValidationError,
			Message: "validation failed",
			Errors:  fields,
		},
	)
}
