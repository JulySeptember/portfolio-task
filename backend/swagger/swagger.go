// swagger/swagger.go

package swagger

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"sync"
)

// =========================
// Embed static assets
// =========================

//go:embed swagger.yml template/*
var StaticAssets embed.FS

type config struct {
	SchemaURL string
	DomID     string
}

var (
	indexTemplate *template.Template
	swaggerYAML   []byte

	loadOnce sync.Once
	loadErr  error
)

// =========================
// init cache
// =========================

func initAssets() {

	loadOnce.Do(func() {

		// =========================
		// template
		// =========================

		t, err := template.ParseFS(
			StaticAssets,
			"template/index.html.tmpl",
		)

		if err != nil {
			loadErr = err
			return
		}

		indexTemplate = t

		// =========================
		// swagger yaml
		// =========================

		data, err := StaticAssets.ReadFile(
			"swagger.yml",
		)

		if err != nil {
			loadErr = err
			return
		}

		swaggerYAML = data
	})
}

// =========================
// Render Swagger UI HTML
// =========================

func renderIndexHTML() ([]byte, error) {

	initAssets()

	if loadErr != nil {
		return nil, loadErr
	}

	c := config{
		SchemaURL: "/api/spec/swagger.yml",
		DomID:     "#root",
	}

	var buf bytes.Buffer

	if err := indexTemplate.Execute(
		&buf,
		c,
	); err != nil {

		return nil, err
	}

	return buf.Bytes(), nil
}

// =========================
// Swagger UI Handler
// =========================

func DocsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		w.WriteHeader(
			http.StatusMethodNotAllowed,
		)

		return
	}

	html, err := renderIndexHTML()

	if err != nil {

		http.Error(
			w,
			"swagger ui unavailable",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	_, _ = w.Write(html)
}

// =========================
// OpenAPI YAML Handler
// =========================

func SpecHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		w.WriteHeader(
			http.StatusMethodNotAllowed,
		)

		return
	}

	initAssets()

	if loadErr != nil {

		http.Error(
			w,
			"swagger spec unavailable",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"text/yaml; charset=utf-8",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	_, _ = w.Write(swaggerYAML)
}
