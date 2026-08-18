package service

import (
	"context"

	"github.com/IBKnight/todo-backend/internal/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUserByUsername(username string) (domain.User, error)
}

type TodoList interface {
	CreateList(ctx context.Context, userId int, list domain.TodoList) (int, error)
	GetUserLists(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetListById(ctx context.Context, listId int, userId int) (domain.TodoList, error)
}

type TodoItem interface {
}
