package todoitem

import (
	"github.com/IBKnight/todo-backend/internal/service"
)

type TodoItemService struct {
	repo service.TodoItemRepository
}

func NewService(repo service.TodoItemRepository) *TodoItemService {
	return &TodoItemService{
		repo: repo,
	}
}
