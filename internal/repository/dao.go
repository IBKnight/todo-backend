package repository

import (
	"github.com/IBKnight/todo-backend/internal/domain"
)

type todoListRow struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

func toRow(l domain.TodoList) todoListRow {
	return todoListRow{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
	}
}

func (r todoListRow) toDomain() domain.TodoList {
	return domain.TodoList{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
	}
}
