package repository

import "github.com/jmoiron/sqlx"

type TodoListRepo struct {
}

func NewTodoListRepo(db *sqlx.DB) *TodoListRepo {
	return &TodoListRepo{}
}
