package todolist

import (
	"context"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/service"
)

type TodoListService struct {
	repo service.TodoListRepository
}

func NewService(repo service.TodoListRepository) *TodoListService {
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

func (s *TodoListService) RemoveList(ctx context.Context, userId int, listId int) error {
	return s.repo.RemoveList(ctx, userId, listId)
}

func (s *TodoListService) UpdateList(ctx context.Context, userId int, list domain.TodoList) (domain.TodoList, error) {
	return s.repo.UpdateList(ctx, userId, list)
}
