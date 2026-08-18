package todoitem

import "github.com/IBKnight/todo-backend/internal/repository"

type TodoItemService struct {
	repo *repository.TodoItemRepo
}

func NewService(repo *repository.TodoItemRepo) *TodoItemService {
	return &TodoItemService{
		repo: repo,
	}
}
