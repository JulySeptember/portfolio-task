package service

import (
	"strings"

	"portfolio/backend/internal/apperr"
	"portfolio/backend/internal/models"
)

// =========================
// NormalizeTaskListQuery
// =========================

func NormalizeTaskListQuery(
	q models.TaskListQuery,
) models.TaskListQuery {

	sortValue := strings.TrimSpace(
		strings.ToLower(string(q.Sort)),
	)

	orderValue := strings.TrimSpace(
		strings.ToUpper(string(q.Order)),
	)

	// =========================
	// defaults
	// =========================

	if sortValue == "" {
		q.Sort = models.TaskSortCreatedAt
	} else {
		q.Sort = models.TaskSort(sortValue)
	}

	if orderValue == "" {
		q.Order = models.TaskOrderDESC
	} else {
		q.Order = models.TaskOrder(orderValue)
	}

	if q.Limit <= 0 {
		q.Limit = 20
	}

	if q.Limit > 100 {
		q.Limit = 100
	}

	if q.Offset < 0 {
		q.Offset = 0
	}

	return q
}

// =========================
// ValidateTaskListQuery
// =========================

func ValidateTaskListQuery(
	q models.TaskListQuery,
) error {

	if !q.Sort.IsValid() {
		return apperr.ErrInvalidSort
	}

	if !q.Order.IsValid() {
		return apperr.ErrInvalidOrder
	}

	if q.Status != "" &&
		!q.Status.IsValid() {

		return apperr.ErrInvalidStatus
	}

	if q.Limit <= 0 ||
		q.Limit > 100 {

		return apperr.ErrInvalidLimit
	}

	if q.Offset < 0 {
		return apperr.ErrInvalidOffset
	}

	return nil
}
