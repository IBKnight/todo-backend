package handler

import "github.com/IBKnight/todo-backend/internal/domain"

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(tokenStr string) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}
