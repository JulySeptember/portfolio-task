// internal/handlers/helpers.go

package handlers

import (
	"net/http"
	"time"

	"portfolio/backend/internal/auth"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/service"
)

// =========================
// parseID
// =========================

func parseID(
	w http.ResponseWriter,
	r *http.Request,
) (int64, bool) {

	id, err := httpx.PathID(r)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			httpx.CodeInvalidID,
			"invalid id",
		)

		return 0, false
	}

	return id, true
}

// =========================
// decodeAndValidate
// =========================

func decodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) bool {

	if err := httpx.DecodeJSON(
		w,
		r,
		dst,
	); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			httpx.CodeInvalidJSON,
			"invalid request body",
		)

		return false
	}

	if errs := httpx.ValidateStruct(
		dst,
	); errs != nil {

		httpx.WriteValidationErrors(
			w,
			errs,
		)

		return false
	}

	return true
}

// =========================
// requireUserID
// =========================

func requireUserID(
	w http.ResponseWriter,
	r *http.Request,
	userSvc *service.UserService,
) (int64, bool) {

	authUser, ok := auth.GetAuthUser(
		r.Context(),
	)

	if !ok || authUser.Sub == "" {

		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			httpx.CodeUnauthorized,
			"missing auth context",
		)

		return 0, false
	}

	user, err := userSvc.GetByAuthUserID(
		r.Context(),
		authUser.Sub,
	)

	if err != nil {

		httpx.HandleError(
			w,
			err,
		)

		return 0, false
	}

	return user.ID, true
}

// =========================
// parseOptionalDueDate
// =========================

func parseOptionalDueDate(
	w http.ResponseWriter,
	value string,
) (*time.Time, bool) {

	t, err := httpx.ParseOptionalTime(
		value,
	)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			httpx.CodeInvalidDueDate,
			"invalid due_date format",
		)

		return nil, false
	}

	return t, true
}
