package todoitem

import (
	"context"

	"github.com/IBKnight/todo-backend/internal/domain"
)

type mockItemRepo struct {
	createFn func(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error)
}

func (m *mockItemRepo) Create(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error) {
	return m.createFn(ctx, userID, listID, item)
}

func (m *mockItemRepo) GetAll(context.Context, int, int) ([]domain.TodoItem, error) {
	return nil, nil
}
func (m *mockItemRepo) GetByID(context.Context, int, int) (domain.TodoItem, error) {
	return domain.TodoItem{}, nil
}
func (m *mockItemRepo) Update(context.Context, int, domain.TodoItem) (domain.TodoItem, error) {
	return domain.TodoItem{}, nil
}
func (m *mockItemRepo) Delete(context.Context, int, int) error {
	return nil
}
