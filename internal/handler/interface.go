package handler

import (
	"context"

	"github.com/IBKnight/todo-backend/internal/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(tokenStr string) (int, error)
}

type TodoList interface {
	CreateList(ctx context.Context, userId int, list domain.TodoList) (int, error)
	GetUserLists(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetListById(ctx context.Context, listId int, userId int) (domain.TodoList, error)
	// UpdateList(userId int, list domain.TodoList) (int, error)
	// RemoveList(userId int, list domain.TodoList) (int, error)
}

type TodoItem interface {
}
