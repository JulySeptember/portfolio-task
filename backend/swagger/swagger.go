package swagger

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ----------------------------
// Embedされた静的ファイル
// ----------------------------

// StaticAssets は swagger.yml とテンプレートを含む
//
//go:embed all:*
var StaticAssets embed.FS

// ----------------------------
// テンプレート用構造体
// ----------------------------
type config struct {
	SchemaURL string
	DomID     string
}

// ----------------------------
// SwaggerYAML は生の swagger.yml を返す
// ----------------------------
func SwaggerYAML() []byte {
	data, err := StaticAssets.ReadFile("swagger.yml")
	if err != nil {
		panic(fmt.Errorf("failed to read swagger.yml: %w", err))
	}
	return data
}

// ----------------------------
// IndexHTML は Swagger UI 用 HTML を返す
// host が空なら相対パス、指定されていれば絶対URLで埋め込む
// ----------------------------
func IndexHTML(host string) []byte {
	t, err := template.ParseFS(StaticAssets, "template/index.html.tmpl")
	if err != nil {
		panic(fmt.Errorf("failed to parse swagger template: %w", err))
	}

	schemaURL := "/api/v1/spec/swagger.yml"
	if host != "" {
		schemaURL = strings.TrimRight(host, "/") + schemaURL
	}

	c := config{
		SchemaURL: schemaURL,
		DomID:     "#root",
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "index.html.tmpl", c); err != nil {
		panic(fmt.Errorf("failed to execute swagger template: %w", err))
	}

	return buf.Bytes()
}

// ----------------------------
// Handler は http.Handler を返す
// ----------------------------
// - /api/v1/spec/swagger.yml  → Swagger定義
// - /api/v1/docs/            → Swagger UI
// ----------------------------
func Handler() http.Handler {
	mux := http.NewServeMux()

	// Swagger YAML 配信
	mux.HandleFunc("/api/v1/spec/swagger.yml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(SwaggerYAML())
	})

	// Swagger UI 配信
	mux.HandleFunc("/api/v1/docs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(IndexHTML(""))
	})

	return mux
}
