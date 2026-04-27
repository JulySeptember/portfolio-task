package handlers

import (
    "reflect"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
    WriteJSON(w, status, map[string]string{"error": msg})
}

func ParseIntOrDefault(s string, def int) int {
    if s == "" {
        return def
    }
    if v, err := strconv.Atoi(s); err == nil {
        return v
    }
    return def
}

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

// DecodeJSON decodes JSON body
func DecodeJSON(r *http.Request, v interface{}) error {
    return json.NewDecoder(r.Body).Decode(v)
}

// SetID sets ID field via reflection
func SetID[T any](t *T, id int64) {
    v := reflect.ValueOf(t).Elem()
    field := v.FieldByName("ID")
    if field.IsValid() && field.CanSet() && field.Kind() == reflect.Int64 {
        field.SetInt(id)
    }
}
