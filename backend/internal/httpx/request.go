// internal/httpx/request.go

package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const maxBodySize = 1 << 20 // 1MB

// =========================
// Errors
// =========================

var (
	ErrInvalidPathID = errors.New("invalid path id")
)

// =========================
// JSON
// =========================

func DecodeJSON(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) error {

	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxBodySize,
	)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {

		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		switch {

		case errors.As(err, &syntaxErr):
			return fmt.Errorf("malformed JSON at position %d", syntaxErr.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")

		case errors.As(err, &typeErr):
			if typeErr.Field != "" {
				return fmt.Errorf("invalid value for field %q", typeErr.Field)
			}
			return errors.New("invalid JSON type")

		case errors.Is(err, io.EOF):
			return errors.New("empty request body")

		case strings.Contains(err.Error(), "http: request body too large"):
			return errors.New("request body too large")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown field %s", field)

		default:
			return errors.New("invalid JSON")
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain only one JSON object")
	}

	return nil
}

// =========================
// Path params
// =========================

func PathID(
	r *http.Request,
) (int64, error) {

	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		return 0, ErrInvalidPathID
	}

	return id, nil
}

// =========================
// Query params
// =========================

func QueryInt(
	r *http.Request,
	key string,
	def int,
	min int,
	max int,
) int {

	v := r.URL.Query().Get(key)

	if v == "" {
		return def
	}

	n, err := strconv.Atoi(v)

	if err != nil {
		return def
	}

	if n < min {
		return min
	}

	if max > 0 && n > max {
		return max
	}

	return n
}

// ★追加：string query
func QueryString(
	r *http.Request,
	key string,
	def string,
) string {

	v := r.URL.Query().Get(key)

	if v == "" {
		return def
	}

	return v
}

// =========================
// Time
// =========================

func ParseOptionalTime(
	value string,
) (*time.Time, error) {

	if value == "" {
		return nil, nil
	}

	t, err := time.Parse(
		time.RFC3339,
		value,
	)

	if err != nil {
		return nil, err
	}

	utc := t.UTC()

	return &utc, nil
}
