package handler

import (
	"context"

	"github.com/IBKnight/todo-backend/internal/domain"
)

type AuthorizationService interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(tokenStr string) (int, error)
}

type TodoListService interface {
	CreateList(ctx context.Context, userId int, list domain.TodoList) (int, error)
	GetUserLists(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetListById(ctx context.Context, listId int, userId int) (domain.TodoList, error)
	UpdateList(ctx context.Context, userId int, updatedTodoList domain.TodoList) (domain.TodoList, error)
	RemoveList(ctx context.Context, userId int, listId int) error
}

type TodoItemService interface {
	Create(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error)
	GetAll(ctx context.Context, userID, listID int) ([]domain.TodoItem, error)
	GetByID(ctx context.Context, userID, itemID int) (domain.TodoItem, error)
	Update(ctx context.Context, userID int, item domain.TodoItem) (domain.TodoItem, error)
	Delete(ctx context.Context, userID, itemID int) error
}
