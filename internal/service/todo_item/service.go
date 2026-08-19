package todoitem

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBKnight/todo-backend/internal/domain"
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

func (s *TodoItemService) Create(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error) {
	if err := validate(item); err != nil {
		return domain.TodoItem{}, err
	}
	return s.repo.Create(ctx, userID, listID, item)
}

func (s *TodoItemService) GetAll(ctx context.Context, userID, listID int) ([]domain.TodoItem, error) {
	return s.repo.GetAll(ctx, userID, listID)
}

func (s *TodoItemService) GetByID(ctx context.Context, userID, itemID int) (domain.TodoItem, error) {
	return s.repo.GetByID(ctx, userID, itemID)
}

func (s *TodoItemService) Update(ctx context.Context, userID int, item domain.TodoItem) (domain.TodoItem, error) {
	if err := validate(item); err != nil {
		return domain.TodoItem{}, err
	}
	return s.repo.Update(ctx, userID, item)
}

func (s *TodoItemService) Delete(ctx context.Context, userID, itemID int) error {
	return s.repo.Delete(ctx, userID, itemID)
}

func validate(item domain.TodoItem) error {
	title := strings.TrimSpace(item.Title)

	if title == "" {
		return fmt.Errorf("%w: title is required", domain.ErrValidation)
	}
	if len([]rune(title)) > 200 {
		return fmt.Errorf("%w: title is too long", domain.ErrValidation)
	}
	if len([]rune(item.Description)) > 1000 {
		return fmt.Errorf("%w: description is too long", domain.ErrValidation)
	}

	return nil
}
