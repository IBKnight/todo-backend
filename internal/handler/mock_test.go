package handler

import (
	"context"
	"testing"

	"github.com/IBKnight/todo-backend/internal/domain"
)

type mockListService struct {
	t *testing.T

	createFn  func(ctx context.Context, userId int, list domain.TodoList) (int, error)
	getAllFn  func(ctx context.Context, userId int) ([]domain.TodoList, error)
	getByIDFn func(ctx context.Context, listId, userId int) (domain.TodoList, error)
	updateFn  func(ctx context.Context, userId int, list domain.TodoList) (domain.TodoList, error)
	removeFn  func(ctx context.Context, userId, listId int) error
}

func (m *mockListService) CreateList(ctx context.Context, userId int, list domain.TodoList) (int, error) {
	if m.createFn == nil {
		m.t.Fatalf("unexpected call to CreateList")
	}
	return m.createFn(ctx, userId, list)
}

func (m *mockListService) GetUserLists(ctx context.Context, userId int) ([]domain.TodoList, error) {
	if m.getAllFn == nil {
		m.t.Fatalf("unexpected call to GetUserLists")
	}
	return m.getAllFn(ctx, userId)
}

func (m *mockListService) GetListById(ctx context.Context, listId, userId int) (domain.TodoList, error) {
	if m.getByIDFn == nil {
		m.t.Fatalf("unexpected call to GetListById")
	}
	return m.getByIDFn(ctx, listId, userId)
}

func (m *mockListService) UpdateList(ctx context.Context, userId int, list domain.TodoList) (domain.TodoList, error) {
	if m.updateFn == nil {
		m.t.Fatalf("unexpected call to UpdateList")
	}
	return m.updateFn(ctx, userId, list)
}

func (m *mockListService) RemoveList(ctx context.Context, userId, listId int) error {
	if m.removeFn == nil {
		m.t.Fatalf("unexpected call to RemoveList")
	}
	return m.removeFn(ctx, userId, listId)
}

type mockItemService struct {
	t *testing.T

	createFn  func(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error)
	getAllFn  func(ctx context.Context, userID, listID int) ([]domain.TodoItem, error)
	getByIDFn func(ctx context.Context, userID, itemID int) (domain.TodoItem, error)
	updateFn  func(ctx context.Context, userID int, item domain.TodoItem) (domain.TodoItem, error)
	deleteFn  func(ctx context.Context, userID, itemID int) error
}

func (m *mockItemService) Create(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error) {
	if m.createFn == nil {
		m.t.Fatalf("unexpected call to Create")
	}
	return m.createFn(ctx, userID, listID, item)
}

func (m *mockItemService) GetAll(ctx context.Context, userID, listID int) ([]domain.TodoItem, error) {
	if m.getAllFn == nil {
		m.t.Fatalf("unexpected call to GetAll")
	}
	return m.getAllFn(ctx, userID, listID)
}

func (m *mockItemService) GetByID(ctx context.Context, userID, itemID int) (domain.TodoItem, error) {
	if m.getByIDFn == nil {
		m.t.Fatalf("unexpected call to GetByID")
	}
	return m.getByIDFn(ctx, userID, itemID)
}

func (m *mockItemService) Update(ctx context.Context, userID int, item domain.TodoItem) (domain.TodoItem, error) {
	if m.updateFn == nil {
		m.t.Fatalf("unexpected call to Update")
	}
	return m.updateFn(ctx, userID, item)
}

func (m *mockItemService) Delete(ctx context.Context, userID, itemID int) error {
	if m.deleteFn == nil {
		m.t.Fatalf("unexpected call to Delete")
	}
	return m.deleteFn(ctx, userID, itemID)
}

type mockAuthService struct {
	t *testing.T

	createUserFn    func(user domain.User) (int, error)
	generateTokenFn func(username, password string) (string, error)
	parseTokenFn    func(tokenStr string) (int, error)
}

func (m *mockAuthService) CreateUser(user domain.User) (int, error) {
	if m.createUserFn == nil {
		m.t.Fatalf("unexpected call to CreateUser")
	}
	return m.createUserFn(user)
}

func (m *mockAuthService) GenerateToken(username, password string) (string, error) {
	if m.generateTokenFn == nil {
		m.t.Fatalf("unexpected call to GenerateToken")
	}
	return m.generateTokenFn(username, password)
}

func (m *mockAuthService) ParseToken(tokenStr string) (int, error) {
	if m.parseTokenFn == nil {
		m.t.Fatalf("unexpected call to ParseToken")
	}
	return m.parseTokenFn(tokenStr)
}
