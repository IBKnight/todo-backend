package todolist

import "github.com/IBKnight/todo-backend/internal/repository"

type TodoListService struct {
	repo *repository.TodoListRepo
}

func NewService(repo *repository.TodoListRepo) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}
