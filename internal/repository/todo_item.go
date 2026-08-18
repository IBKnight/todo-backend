package repository

import "github.com/jmoiron/sqlx"

type TodoItemRepo struct {
	db *sqlx.DB
}

func NewTodoItemRepo(db *sqlx.DB) *TodoItemRepo {
	return &TodoItemRepo{db}
}
