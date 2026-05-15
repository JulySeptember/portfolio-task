package models

import "strings"

type TaskListQuery struct {
	Status string
	Sort   string
	Order  string
	Limit  int
	Offset int
}

func (q TaskListQuery) Normalize() TaskListQuery {
	q.Status = strings.TrimSpace(q.Status)

	q.Sort = normalizeSort(q.Sort)
	q.Order = normalizeOrder(q.Order)

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

func normalizeSort(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "due_date":
		return "due_date"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}

func normalizeOrder(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "asc":
		return "ASC"
	default:
		return "DESC"
	}
}
