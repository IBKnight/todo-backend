package todolist

import "github.com/IBKnight/todo-backend/internal/repository"

type TodoList struct {
	repo *repository.TodoListRepo
}

func NewService(repo *repository.TodoListRepo) *TodoList {
	return &TodoList{
		repo: repo,
	}
}
