// internal/models/task_query.go

package models

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
