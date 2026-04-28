package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// WriteJSON writes v as JSON with given status.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes a JSON error message.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// ParseIntOrDefault parses s as int, returns def on error or empty string.
func ParseIntOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}

// ExtractIDFromPath extracts a single numeric ID from path after prefix.
func ExtractIDFromPath(path, prefix string) (int64, error) {
	if !strings.HasPrefix(path, prefix) {
		return 0, strconv.ErrSyntax
	}
	rest := strings.TrimPrefix(path, prefix)
	rest = strings.Trim(rest, "/")
	if rest == "" || strings.Contains(rest, "/") {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(rest, 10, 64)
}

// DecodeJSON decodes JSON body into v with safety measures:
// - limits body size to 1MB
// - disallows unknown fields
// Returns a descriptive error for client responses.
func DecodeJSON(r *http.Request, v interface{}) error {
	// Limit request body size to prevent DoS via large payloads.
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1MB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		// Provide clearer error messages for common cases.
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			return errors.New("malformed JSON")
		case errors.Is(err, http.ErrBodyReadAfterClose):
			return errors.New("request body closed")
		case errors.As(err, &unmarshalTypeErr):
			return errors.New("invalid value for field " + unmarshalTypeErr.Field)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			// err.Error() looks like: "json: unknown field \"foo\""
			return errors.New("unknown field in JSON")
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")
		case err.Error() == "http: request body too large":
			return errors.New("request body too large")
		default:
			return err
		}
	}

	// Ensure there is no extra data after the JSON object.
	if dec.More() {
		return errors.New("multiple JSON objects in body")
	}

	return nil
}

// SetID sets the ID field on a pointer-to-struct safely via reflection.
func SetID[T any](t *T, id int64) {
	if t == nil {
		return
	}
	v := reflect.ValueOf(t)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	e := v.Elem()
	if !e.IsValid() || e.Kind() != reflect.Struct {
		return
	}
	field := e.FieldByName("ID")
	if !field.IsValid() || !field.CanSet() {
		return
	}
	switch field.Kind() {
	case reflect.Int64:
		field.SetInt(id)
	case reflect.Int, reflect.Int32:
		field.SetInt(id)
	default:
	}
}
