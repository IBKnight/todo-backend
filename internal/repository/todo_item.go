package repository

import "github.com/jmoiron/sqlx"

type TodoItemRepo struct {
}

func NewTodoItemRepo(db *sqlx.DB) *TodoItemRepo {
	return &TodoItemRepo{}
}
