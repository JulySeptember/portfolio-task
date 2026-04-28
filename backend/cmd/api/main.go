package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"portfolio/backend/internal/config"
	"portfolio/backend/internal/repository"
	"portfolio/backend/internal/router"
	"portfolio/backend/internal/service"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// globalDB は Lambda のコールドスタート時または runLocal() で初期化される *sql.DB
	globalDB *sql.DB

	// mux は Lambda 用のアダプタやローカルサーバで使う http.Handler
	mux http.Handler

	// docsTmpl は Swagger UI 用テンプレートをキャッシュするための変数
	docsTmpl *template.Template
)

// init はパッケージ初期化時に呼ばれる（Lambda のコールドスタートで実行される想定）
// RUN_MODE=local のときはローカル起動処理を main() 側で行うためここでは何もしない。
func init() {
	if os.Getenv("RUN_MODE") == "local" {
		// ローカル実行では init で DB を初期化しない（runLocal() で行う）
		return
	}

	// Lambda 実行環境（local 以外）ではここで DB を初期化しておく
	// これによりコールドスタート時に DB 接続を確立できる
	if err := initDB(); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	// Swagger テンプレートをコールドスタート時に一度だけパースしてキャッシュする
	if t, err := template.ParseFiles("./swagger/template/index.html.tmpl"); err != nil {
		// テンプレートが無ければログに出すが起動は続行（serveDocs はエラーを返す）
		log.Printf("warning: failed to parse docs template: %v", err)
	} else {
		docsTmpl = t
	}

	// ルーターを構築して正規化ミドルウェア＋ロギングミドルウェアでラップする
	// Lambda ではこの mux を httpadapter に渡してハンドリングする
	mux = loggingMiddleware(router.NormalizePathMiddleware(buildMux()))
}

func main() {
	// ローカル実行モードなら runLocal() を呼んで通常の HTTP サーバを起動する
	if os.Getenv("RUN_MODE") == "local" {
		runLocal()
		return
	}

	// Lambda 実行モード：既に init() で構築した mux をアダプタに渡して起動
	adapter := httpadapter.New(mux)
	lambda.Start(adapter.ProxyWithContext)
}

// initDB は globalDB を初期化する（ConnectDBFromEnv を呼ぶだけに集約）
// Lambda の init() から呼ばれる想定
func initDB() error {
	if globalDB != nil {
		// 既に初期化済みなら何もしない
		return nil
	}
	db, err := config.ConnectDBFromEnv()
	if err != nil {
		return err
	}
	globalDB = db
	return nil
}

// CloseDB は globalDB を閉じるユーティリティ（テストやシャットダウンで使える）
func CloseDB() error {
	if globalDB == nil {
		return nil
	}
	err := globalDB.Close()
	globalDB = nil
	return err
}

// runLocal はローカル開発用のエントリポイント。
// - DB を ConnectDBFromEnv で初期化
// - リポジトリ/サービスを作成
// - *http.ServeMux にルートを登録
// - Swagger の静的配信を設定
// - 標準的な http.Server を起動（ログミドルウェアを適用）
func runLocal() {
	// DB 接続（環境変数から DSN 等を読む）
	db, err := config.ConnectDBFromEnv()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	// ローカルではプロセス終了時に DB を閉じる
	defer func() { _ = db.Close() }()

	// globalDB にセットして buildMux と同じ形でサービスを初期化できるようにする
	globalDB = db

	// リポジトリとサービスを作成（依存注入）
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	taskRepo := repository.NewTaskRepo(db)
	taskSvc := service.NewTaskService(taskRepo)

	// 具体的な *http.ServeMux を作成してルートを登録する
	localMux := http.NewServeMux()
	router.RegisterRoutes(localMux, userSvc, taskSvc)

	// Swagger UI と OpenAPI スキーマの静的配信を設定
	localMux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))
	localMux.HandleFunc("/api/v1/docs", serveDocs)

	// ローカルでもテンプレートを一度だけパースしてキャッシュする（起動時）
	if docsTmpl == nil {
		if t, err := template.ParseFiles("./swagger/template/index.html.tmpl"); err != nil {
			log.Printf("warning: failed to parse docs template: %v", err)
		} else {
			docsTmpl = t
		}
	}

	// サーバ設定（タイムアウト等）
	addr := ":" + getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(router.NormalizePathMiddleware(localMux)), // 正規化→ログの順でラップ
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("local server listening on %s (RUN_MODE=%s)", addr, os.Getenv("RUN_MODE"))
	// ListenAndServe はブロッキング。エラーがあればログ出力して終了する
	log.Fatal(srv.ListenAndServe())
}

// buildMux は Lambda 用に返す *http.ServeMux を構築する関数。
// init() から呼ばれて mux を作る。globalDB を使ってサービスを初期化する点に注意。
func buildMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Swagger とドキュメント配信
	mux.Handle("/api/v1/spec/", http.StripPrefix("/api/v1/spec/", http.FileServer(http.Dir("./swagger"))))
	mux.HandleFunc("/api/v1/docs", serveDocs)

	// globalDB を使ってリポジトリ/サービスを初期化
	userRepo := repository.NewUserRepository(globalDB)
	userSvc := service.NewUserService(userRepo)

	taskRepo := repository.NewTaskRepo(globalDB)
	taskSvc := service.NewTaskService(taskRepo)

	// ルートを登録（router.RegisterRoutes は具体的なハンドラを mux に登録する）
	router.RegisterRoutes(mux, userSvc, taskSvc)

	return mux
}

// serveDocs は Swagger UI のテンプレートを読み込み、HTML を返すハンドラ。
// テンプレート内で SchemaURL を参照して Swagger UI を表示する。
func serveDocs(w http.ResponseWriter, r *http.Request) {
	// 既にキャッシュされたテンプレートがあればそれを使う
	if docsTmpl != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := docsTmpl.Execute(w, map[string]string{
			"SchemaURL": "/api/v1/spec/swagger.yml",
			"DomID":     "#root",
		}); err != nil {
			http.Error(w, "docs not available", http.StatusInternalServerError)
		}
		return
	}

	// キャッシュが無ければフォールバックで都度パース（起動時に一度だけパースすることを推奨）
	tmpl, err := template.ParseFiles("./swagger/template/index.html.tmpl")
	if err != nil {
		http.Error(w, "docs not available", http.StatusInternalServerError)
		return
	}
	data := map[string]string{
		"SchemaURL": "/api/v1/spec/swagger.yml",
		"DomID":     "#root",
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, data)
}

// getEnv は環境変数を読み、未設定ならデフォルトを返すユーティリティ
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// loggingMiddleware はリクエスト開始/終了をログ出力する簡易ミドルウェア。
// - リモートアドレスを X-Forwarded-For から優先して取得
// - 開始ログと完了ログに経過時間を出力
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// リモートアドレスの取得（X-Forwarded-For を優先）
		remote := r.RemoteAddr
		if remote == "" {
			remote = "-"
		}
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			remote = strings.TrimSpace(parts[0])
		}

		// リクエスト開始ログ
		log.Printf("started %s %s from %s", r.Method, r.URL.Path, remote)

		// 実際のハンドラを呼び出す
		next.ServeHTTP(w, r)

		// 完了ログ（処理時間を表示）
		log.Printf("completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
