package swagger

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

// config はテンプレートに渡すパラメータを定義します。
type config struct {
	SchemaURL string
	DomID     string
}

// StaticAssets はディレクトリ全体を保持します（http.FileServer 用）。
// ディレクティブの前にスペースを入れず、説明コメントは上の行に配置します。
//
//go:embed all:*
var StaticAssets embed.FS

// SwaggerYAML は埋め込まれた swagger.yml の生データを返します。
func SwaggerYAML() []byte {
	// ファイル名は embed したルートからの相対パスです
	data, err := StaticAssets.ReadFile("swagger.yml")
	if err != nil {
		// 初期化時にファイルがない場合は致命的なため panic させています
		panic(fmt.Errorf("failed to read swagger.yml: %w", err))
	}
	return data
}

// IndexHTML はテンプレートをパースし、指定されたホスト情報を埋め込んだ HTML を返します。
// host が空の場合は相対パスを使用します。
func IndexHTML(host string) []byte {
	// StaticAssets からテンプレートファイルを読み込みます
	// 第2引数は embed されたルートからの相対パスを指定します
	t, err := template.ParseFS(StaticAssets, "template/index.html.tmpl")
	if err != nil {
		panic(fmt.Errorf("failed to parse swagger template: %w", err))
	}

	schemaURL := "/api/v1/spec/swagger.yml"
	if host != "" {
		schemaURL = fmt.Sprintf("%s%s", host, schemaURL)
	}

	c := config{
		SchemaURL: schemaURL,
		DomID:     "#root",
	}

	var buffer bytes.Buffer
	// テンプレート実行時はファイル名（ベース名）を指定します
	if err := t.ExecuteTemplate(&buffer, "index.html.tmpl", c); err != nil {
		panic(fmt.Errorf("failed to execute swagger template: %w", err))
	}

	return buffer.Bytes()
}
