package swagger

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed swagger.yml template/*
var StaticAssets embed.FS

type config struct {
	SchemaURL string
	DomID     string
}

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

func IndexHTML() []byte {

	t, err := template.ParseFS(
		StaticAssets,
		"template/index.html.tmpl",
	)
	if err != nil {
		panic(
			fmt.Errorf(
				"failed to parse template: %w",
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
				"failed to execute template: %w",
				err,
			),
		)
	}

	return buf.Bytes()
}

func DocsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.Write(IndexHTML())
}

func SpecHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/yaml",
	)

	w.Write(SwaggerYAML())
}
