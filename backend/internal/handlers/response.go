package handlers

import (
	"encoding/json"
	"net/http"
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

	_ = json.NewEncoder(w).Encode(v)
}

// =========================
// error response
// =========================

func WriteError(
	w http.ResponseWriter,
	status int,
	msg string,
) {
	WriteJSON(w, status, map[string]string{
		"error": msg,
	})
}

// =========================
// validation error response
// =========================

func WriteValidationErrors(
	w http.ResponseWriter,
	errs map[string]string,
) {
	WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
		"errors": errs,
	})
}
