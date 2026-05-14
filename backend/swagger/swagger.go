package swagger

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

// =========================
// Embed static assets
// =========================
//
// swagger.yml
// template/index.html.tmpl
//
// をバイナリへ埋め込む
// Lambda deploy 時も単一バイナリで完結
// =========================

//go:embed swagger.yml template/*
var StaticAssets embed.FS

type config struct {
	SchemaURL string
	DomID     string
}

// =========================
// Load swagger.yml
// =========================

func SwaggerYAML() []byte {

	data, err := StaticAssets.ReadFile(
		"swagger.yml",
	)

	if err != nil {

		panic(
			fmt.Errorf(
				"failed to read swagger.yml: %w",
				err,
			),
		)
	}

	return data
}

// =========================
// Render Swagger UI HTML
// =========================

func IndexHTML() []byte {

	t, err := template.ParseFS(
		StaticAssets,
		"template/index.html.tmpl",
	)

	if err != nil {

		panic(
			fmt.Errorf(
				"failed to parse swagger template: %w",
				err,
			),
		)
	}

	c := config{
		SchemaURL: "/api/v1/spec/swagger.yml",
		DomID:     "#root",
	}

	var buf bytes.Buffer

	if err := t.Execute(
		&buf,
		c,
	); err != nil {

		panic(
			fmt.Errorf(
				"failed to execute swagger template: %w",
				err,
			),
		)
	}

	return buf.Bytes()
}

// =========================
// Swagger UI Handler
// =========================

func DocsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(
		IndexHTML(),
	)
}

// =========================
// OpenAPI YAML Handler
// =========================

func SpecHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/yaml",
	)

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(
		SwaggerYAML(),
	)
}
