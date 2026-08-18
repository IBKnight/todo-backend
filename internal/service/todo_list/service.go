package todolist

import (
	"context"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/repository"
)

type TodoListService struct {
	repo *repository.TodoListRepo
}

func NewService(repo *repository.TodoListRepo) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (s *TodoListService) CreateList(ctx context.Context, userId int, list domain.TodoList) (int, error) {
	return s.repo.CreateList(ctx, userId, list)
}

func (s *TodoListService) GetUserLists(ctx context.Context, userId int) ([]domain.TodoList, error) {
	return s.repo.GetUserLists(ctx, userId)

}

func (s *TodoListService) GetListById(ctx context.Context, listId int, userId int) (domain.TodoList, error) {
	return s.repo.GetListById(ctx, listId, userId)
}
