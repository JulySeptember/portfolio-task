// internal/models/task_query.go

package models

import (
	"strings"

	"portfolio/backend/internal/apperr"
)

// =========================
// sort enum
// =========================

type TaskSort string

const (
	TaskSortCreatedAt TaskSort = "created_at"
	TaskSortDueDate   TaskSort = "due_date"
)

func (s TaskSort) IsValid() bool {

	switch s {

	case TaskSortCreatedAt,
		TaskSortDueDate:

		return true

	default:
		return false
	}
}

// =========================
// SQL column
// =========================

func (s TaskSort) Column() string {

	switch s {

	case TaskSortDueDate:
		return "due_date"

	default:
		return "created_at"
	}
}

// =========================
// order enum
// =========================

type TaskOrder string

const (
	TaskOrderASC  TaskOrder = "ASC"
	TaskOrderDESC TaskOrder = "DESC"
)

func (o TaskOrder) IsValid() bool {

	switch o {

	case TaskOrderASC,
		TaskOrderDESC:

		return true

	default:
		return false
	}
}

// =========================
// SQL order
// =========================

func (o TaskOrder) SQL() string {

	switch o {

	case TaskOrderASC:
		return "ASC"

	default:
		return "DESC"
	}
}

// =========================
// query
// =========================

type TaskListQuery struct {
	Status TaskStatus

	Sort  TaskSort
	Order TaskOrder

	Limit  int
	Offset int
}

// =========================
// list result
// =========================

type TaskListResult struct {
	Items []Task `json:"items"`
	Total int64  `json:"total"`
}

// =========================
// normalize
// =========================

func (q *TaskListQuery) Normalize() {

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
		q.Sort = TaskSortCreatedAt
	} else {
		q.Sort = TaskSort(sortValue)
	}

	if orderValue == "" {
		q.Order = TaskOrderDESC
	} else {
		q.Order = TaskOrder(orderValue)
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
}

// =========================
// validate
// =========================

func (q TaskListQuery) Validate() error {

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
