package handlers

import (
	"errors"
	"net/http"

	"portfolio/backend/internal/dto"
	"portfolio/backend/internal/models"
	"portfolio/backend/internal/service"
)

func DecodeCreateUser(w http.ResponseWriter, r *http.Request) (*models.User, int, error) {
	var req dto.CreateUserRequest
	if err := DecodeJSON(w, r, &req); err != nil {
		return nil, http.StatusBadRequest, errors.New("invalid body")
	}
	if req.Email == "" {
		return nil, http.StatusBadRequest, errors.New("email is required")
	}

	u := &models.User{
		Email:       req.Email,
		DisplayName: req.DisplayName,
	}
	return u, 0, nil
}

func MergeUpdateUser(w http.ResponseWriter, r *http.Request, existing *models.User) (int, error) {
	var req dto.UpdateUserRequest
	if err := DecodeJSON(w, r, &req); err != nil {
		return http.StatusBadRequest, errors.New("invalid body")
	}
	if req.DisplayName != "" {
		existing.DisplayName = req.DisplayName
	}
	return 0, nil
}

func NewUserHandler(svc *service.UserService) *BaseHandler[models.User] {
	return NewBaseHandlerWithDTO(svc, DecodeCreateUser, MergeUpdateUser)
}

type UserHandlerWrapper struct {
	base *BaseHandler[models.User]
}

func NewUserHandlerWrapper(svc *service.UserService) *UserHandlerWrapper {
	return &UserHandlerWrapper{base: NewUserHandler(svc)}
}

func (w *UserHandlerWrapper) Create(rw http.ResponseWriter, r *http.Request) {
	w.base.Create(rw, r)
}
func (w *UserHandlerWrapper) List(rw http.ResponseWriter, r *http.Request) {
	w.base.List(rw, r)
}
func (w *UserHandlerWrapper) Get(rw http.ResponseWriter, r *http.Request, id int64) {
	w.base.Get(rw, r, id)
}
func (w *UserHandlerWrapper) HandleUpdate(rw http.ResponseWriter, r *http.Request, id int64) {
	w.base.HandleUpdate(rw, r, id)
}
func (w *UserHandlerWrapper) Delete(rw http.ResponseWriter, r *http.Request, id int64) {
	w.base.Delete(rw, r, id)
}
