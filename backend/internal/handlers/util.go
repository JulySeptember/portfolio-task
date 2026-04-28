package handlers

import (
        "encoding/json"
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

// DecodeJSON decodes JSON body into v.
func DecodeJSON(r *http.Request, v interface{}) error {
        return json.NewDecoder(r.Body).Decode(v)
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
