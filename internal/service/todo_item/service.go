package todoitem

import "github.com/IBKnight/todo-backend/internal/repository"

type TodoItem struct {
	repo *repository.TodoItemRepo
}

func NewService(repo *repository.TodoItemRepo) *TodoItem {
	return &TodoItem{
		repo: repo,
	}
}
