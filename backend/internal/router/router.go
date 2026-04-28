package router

import (
	"net/http"
	"strings"

	"portfolio/backend/internal/handlers"
	"portfolio/backend/internal/service"
)

// NormalizePathMiddleware removes trailing slash from URL.Path except when path == "/".
// It mutates r.URL.Path for downstream handlers (no redirect).
func NormalizePathMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := r.URL.Path; len(p) > 1 && strings.HasSuffix(p, "/") {
			r.URL.Path = strings.TrimRight(p, "/")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
		}
		next.ServeHTTP(w, r)
	})
}

// ResourceHandler はルーターが期待するハンドラの振る舞いを定義します。
// TaskHandlerWrapper / UserHandlerWrapper はこのインターフェースを満たします。
type ResourceHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id int64)
	HandleUpdate(w http.ResponseWriter, r *http.Request, id int64)
	Delete(w http.ResponseWriter, r *http.Request, id int64)
}

// BaseRouter は /prefix と /prefix/{id} の基本ルーティングを提供します。
type BaseRouter struct {
	handler ResourceHandler
	prefix  string
}

func NewBaseRouter(prefix string, h ResourceHandler) *BaseRouter {
	return &BaseRouter{handler: h, prefix: prefix}
}

func (r *BaseRouter) Register(mux *http.ServeMux) {
	// collection path (exact)
	mux.HandleFunc(r.prefix, func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.handler.List(w, req)
		case http.MethodPost:
			r.handler.Create(w, req)
		default:
			handlers.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// item path (with id)
	itemPrefix := r.prefix + "/"
	mux.HandleFunc(itemPrefix, func(w http.ResponseWriter, req *http.Request) {
		id, err := handlers.ExtractIDFromPath(req.URL.Path, itemPrefix)
		if err != nil || id <= 0 {
			handlers.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}

		switch req.Method {
		case http.MethodGet:
			r.handler.Get(w, req, id)
		case http.MethodPut, http.MethodPatch:
			r.handler.HandleUpdate(w, req, id)
		case http.MethodDelete:
			r.handler.Delete(w, req, id)
		default:
			handlers.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})
}

// RegisterRoutes はアプリケーションのユーザーとタスクのルートを一括登録します.
// mux: サーバの ServeMux
// userSvc: UserService のインスタンス
// taskSvc: TaskService のインスタンス
func RegisterRoutes(mux *http.ServeMux, userSvc *service.UserService, taskSvc *service.TaskService) {
	// Users
	userHandler := handlers.NewUserHandlerWrapper(userSvc)
	userRouter := NewBaseRouter("/api/v1/users", userHandler)
	userRouter.Register(mux)

	// Tasks
	taskHandler := handlers.NewTaskHandlerWrapper(taskSvc)
	taskRouter := NewBaseRouter("/api/v1/tasks", taskHandler)
	taskRouter.Register(mux)
}
