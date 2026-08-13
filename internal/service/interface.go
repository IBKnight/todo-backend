package service

import "github.com/IBKnight/todo-backend/internal/domain"

type Authorization interface {
	CreateUser(user domain.User) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}
