package handlers

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

// =========================
// query helper
// =========================

func ParseIntOrDefault(
	s string,
	def int,
) int {

	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}

	return v
}

// =========================
// ID helper
// =========================

// ParseID extracts ID from URL path.
// example:
// /users/123
func ParseID(
	r *http.Request,
	resource string,
) (int64, bool) {

	parts := strings.Split(
		strings.Trim(r.URL.Path, "/"),
		"/",
	)

	// expected:
	// api/v1/users/123
	if len(parts) != 4 {
		return 0, false
	}

	if parts[0] != "api" ||
		parts[1] != "v1" ||
		parts[2] != resource {

		return 0, false
	}

	id, err := strconv.ParseInt(
		parts[3],
		10,
		64,
	)

	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

// =========================
// optional time parser
// =========================

// ParseOptionalTime parses RFC3339 datetime.
// empty string returns nil.
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

	return &t, nil
}

// =========================
// JSON helper
// =========================

// DecodeJSON safely decodes JSON request body.
//
// features:
// - reject unknown fields
// - limit body size
// - detailed syntax errors
func DecodeJSON(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) error {

	const maxBodySize = 1 << 20 // 1MB

	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxBodySize,
	)

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {

		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {

		case errors.As(err, &syntaxErr):

			return fmt.Errorf(
				"malformed JSON at position %d",
				syntaxErr.Offset,
			)

		case errors.Is(err, io.ErrUnexpectedEOF):

			return fmt.Errorf("malformed JSON")

		case errors.As(err, &unmarshalTypeErr):

			return fmt.Errorf(
				"invalid value for field %q at position %d",
				unmarshalTypeErr.Field,
				unmarshalTypeErr.Offset,
			)

		case errors.Is(err, io.EOF):

			return fmt.Errorf("empty request body")

		case strings.HasPrefix(
			err.Error(),
			"http: request body too large",
		):

			return fmt.Errorf(
				"request body too large, max %d bytes",
				maxBodySize,
			)

		default:

			return fmt.Errorf(
				"invalid JSON: %v",
				err,
			)
		}
	}

	// prevent multiple JSON objects
	if err := dec.Decode(&struct{}{}); err != io.EOF {

		return fmt.Errorf(
			"body must contain only one JSON object",
		)
	}

	return nil
}
