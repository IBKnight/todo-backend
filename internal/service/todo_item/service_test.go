package todoitem

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/IBKnight/todo-backend/internal/domain"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		item    domain.TodoItem
		wantErr error
	}{
		{
			name:    "valid item",
			item:    domain.TodoItem{Title: "buy milk", Description: "2%"},
			wantErr: nil,
		},
		{
			name:    "empty title",
			item:    domain.TodoItem{Title: ""},
			wantErr: domain.ErrValidation,
		},
		{
			name:    "whitespace only title",
			item:    domain.TodoItem{Title: "   "},
			wantErr: domain.ErrValidation,
		},
		{
			name:    "title too long",
			item:    domain.TodoItem{Title: strings.Repeat("a", 201)},
			wantErr: domain.ErrValidation,
		},
		{
			name:    "description too long",
			item:    domain.TodoItem{Title: "ok", Description: strings.Repeat("a", 1001)},
			wantErr: domain.ErrValidation,
		},
		{
			name:    "cyrillic title at limit",
			item:    domain.TodoItem{Title: strings.Repeat("я", 200)},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.item)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("actual: %v, expected: %v", err, tt.wantErr)

			}
		})
	}
}

func TestTodoItemService_Create(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.TodoItem
		repoErr  error
		wantErr  error
		wantCall bool
	}{
		{
			name:     "success",
			input:    domain.TodoItem{Title: "buy milk"},
			wantCall: true,
		},
		{
			name:    "empty title",
			input:   domain.TodoItem{Title: ""},
			wantErr: domain.ErrValidation,
		},
		{
			name:     "list not found",
			input:    domain.TodoItem{Title: "buy milk"},
			repoErr:  domain.ErrListNotFound,
			wantErr:  domain.ErrListNotFound,
			wantCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			repo := &mockItemRepo{
				createFn: func(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error) {
					called = true
					if tt.repoErr != nil {
						return domain.TodoItem{}, tt.repoErr
					}
					item.ID = 1
					return item, nil
				},
			}

			svc := NewService(repo)
			_, err := svc.Create(context.Background(), 1, 10, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if called != tt.wantCall {
				t.Errorf("repo called = %v, want %v", called, tt.wantCall)
			}
		})
	}
}
