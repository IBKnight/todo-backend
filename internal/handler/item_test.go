package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/dto"
)

const testListID = 10

func TestHandler_createItem(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		body       string
		svcErr     error
		wantStatus int
		wantCall   bool
	}{
		{
			name:       "success",
			param:      "10",
			body:       `{"title":"buy milk","description":"2%"}`,
			wantStatus: http.StatusCreated,
			wantCall:   true,
		},
		{
			name:       "invalid list id",
			param:      "abc",
			body:       `{"title":"buy milk"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty title rejected by binding",
			param:      "10",
			body:       `{"title":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed json",
			param:      "10",
			body:       `{"title":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "list not found",
			param:      "999",
			body:       `{"title":"buy milk"}`,
			svcErr:     domain.ErrListNotFound,
			wantStatus: http.StatusNotFound,
			wantCall:   true,
		},
		{
			name:       "validation error from service",
			param:      "10",
			body:       `{"title":"buy milk"}`,
			svcErr:     fmt.Errorf("%w: title is required", domain.ErrValidation),
			wantStatus: http.StatusBadRequest,
			wantCall:   true,
		},
		{
			name:       "service error",
			param:      "10",
			body:       `{"title":"buy milk"}`,
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUserID, gotListID int
			var gotItem domain.TodoItem

			svc := &mockItemService{
				t: t,
				createFn: func(ctx context.Context, uID, lID int, item domain.TodoItem) (domain.TodoItem, error) {
					called = true
					gotUserID, gotListID, gotItem = uID, lID, item
					if tt.svcErr != nil {
						return domain.TodoItem{}, tt.svcErr
					}
					item.ID = 1
					item.ListID = lID
					return item, nil
				},
			}

			h := &Handler{item: svc}
			r := setupRouter(t, http.MethodPost, "/lists/:id/items", testUserID, h.createItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/lists/"+tt.param+"/items",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if called != tt.wantCall {
				t.Fatalf("service called = %v, want %v", called, tt.wantCall)
			}
			if !called {
				return
			}

			if gotUserID != testUserID {
				t.Errorf("service got userID %d, want %d", gotUserID, testUserID)
			}
			if tt.param == "10" && gotListID != testListID {
				t.Errorf("service got listID %d, want %d", gotListID, testListID)
			}
			if tt.svcErr == nil && gotItem.Title != "buy milk" {
				t.Errorf("service got title %q, want %q", gotItem.Title, "buy milk")
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			if loc := w.Header().Get("Location"); loc != "/api/items/1" {
				t.Errorf("got Location %q, want /api/items/1", loc)
			}

			var resp dto.ItemResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.ID != 1 || resp.Title != "buy milk" {
				t.Errorf("got %+v", resp)
			}
		})
	}
}

func TestHandler_getAllItems(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		want := []domain.TodoItem{
			{ID: 1, ListID: testListID, Title: "buy milk", Done: false},
			{ID: 2, ListID: testListID, Title: "walk dog", Done: true},
		}

		var gotUserID, gotListID int
		svc := &mockItemService{
			t: t,
			getAllFn: func(ctx context.Context, uID, lID int) ([]domain.TodoItem, error) {
				gotUserID, gotListID = uID, lID
				return want, nil
			},
		}

		h := &Handler{item: svc}
		r := setupRouter(t, http.MethodGet, "/lists/:id/items", testUserID, h.getAllItems)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists/10/items", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("got status %d, want 200, body: %s", w.Code, w.Body.String())
		}
		if gotUserID != testUserID || gotListID != testListID {
			t.Errorf("service got (%d, %d), want (%d, %d)", gotUserID, gotListID, testUserID, testListID)
		}

		var resp dto.ItemsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(resp.Data) != 2 {
			t.Fatalf("got %d items, want 2", len(resp.Data))
		}
		if !resp.Data[1].Done {
			t.Error("second item should be done")
		}
	})

	t.Run("empty list serializes as array", func(t *testing.T) {
		svc := &mockItemService{
			t: t,
			getAllFn: func(context.Context, int, int) ([]domain.TodoItem, error) {
				return nil, nil
			},
		}

		h := &Handler{item: svc}
		r := setupRouter(t, http.MethodGet, "/lists/:id/items", testUserID, h.getAllItems)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists/10/items", nil))

		if strings.Contains(w.Body.String(), "null") {
			t.Errorf("expected empty array, got %s", w.Body.String())
		}
	})

	t.Run("invalid list id", func(t *testing.T) {
		svc := &mockItemService{t: t}

		h := &Handler{item: svc}
		r := setupRouter(t, http.MethodGet, "/lists/:id/items", testUserID, h.getAllItems)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists/abc/items", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("got status %d, want 400", w.Code)
		}
	})

	t.Run("no user id in context", func(t *testing.T) {
		svc := &mockItemService{t: t}

		h := &Handler{item: svc}
		r := setupRouter(t, http.MethodGet, "/lists/:id/items", 0, h.getAllItems)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists/10/items", nil))

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("got status %d, want 500", w.Code)
		}
	})
}

func TestHandler_getItemById(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		svcItem    domain.TodoItem
		svcErr     error
		wantStatus int
		wantCall   bool
	}{
		{
			name:       "success",
			param:      "42",
			svcItem:    domain.TodoItem{ID: 42, ListID: testListID, Title: "buy milk", Done: true},
			wantStatus: http.StatusOK,
			wantCall:   true,
		},
		{
			name:       "not found",
			param:      "999",
			svcErr:     domain.ErrItemNotFound,
			wantStatus: http.StatusNotFound,
			wantCall:   true,
		},
		{
			name:       "invalid item id",
			param:      "abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			param:      "42",
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUserID, gotItemID int

			svc := &mockItemService{
				t: t,
				getByIDFn: func(ctx context.Context, uID, iID int) (domain.TodoItem, error) {
					called = true
					gotUserID, gotItemID = uID, iID
					if tt.svcErr != nil {
						return domain.TodoItem{}, tt.svcErr
					}
					return tt.svcItem, nil
				},
			}

			h := &Handler{item: svc}
			r := setupRouter(t, http.MethodGet, "/lists/:id/items/:item_id", testUserID, h.getItemById)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists/10/items/"+tt.param, nil))

			if w.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if called != tt.wantCall {
				t.Fatalf("service called = %v, want %v", called, tt.wantCall)
			}
			if !called {
				return
			}

			if gotUserID != testUserID {
				t.Errorf("service got userID %d, want %d", gotUserID, testUserID)
			}
			if tt.param == "42" && gotItemID != 42 {
				t.Errorf("service got itemID %d, want 42", gotItemID)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp dto.ItemResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.ID != tt.svcItem.ID || resp.Title != tt.svcItem.Title || resp.Done != tt.svcItem.Done {
				t.Errorf("got %+v, want %+v", resp, tt.svcItem)
			}
		})
	}
}

func TestHandler_updateItem(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		body       string
		svcErr     error
		wantStatus int
		wantCall   bool
		wantDone   bool
	}{
		{
			name:       "success mark done",
			param:      "42",
			body:       `{"title":"buy milk","description":"2%","done":true}`,
			wantStatus: http.StatusOK,
			wantCall:   true,
			wantDone:   true,
		},
		{
			name:       "success mark undone",
			param:      "42",
			body:       `{"title":"buy milk","done":false}`,
			wantStatus: http.StatusOK,
			wantCall:   true,
			wantDone:   false,
		},
		{
			name:       "invalid item id",
			param:      "abc",
			body:       `{"title":"buy milk"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty title rejected by binding",
			param:      "42",
			body:       `{"title":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed json",
			param:      "42",
			body:       `{"title":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			param:      "999",
			body:       `{"title":"buy milk"}`,
			svcErr:     domain.ErrItemNotFound,
			wantStatus: http.StatusNotFound,
			wantCall:   true,
		},
		{
			name:       "service error",
			param:      "42",
			body:       `{"title":"buy milk"}`,
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUserID int
			var gotItem domain.TodoItem

			svc := &mockItemService{
				t: t,
				updateFn: func(ctx context.Context, uID int, item domain.TodoItem) (domain.TodoItem, error) {
					called = true
					gotUserID, gotItem = uID, item
					if tt.svcErr != nil {
						return domain.TodoItem{}, tt.svcErr
					}
					item.ListID = testListID
					return item, nil
				},
			}

			h := &Handler{item: svc}
			r := setupRouter(t, http.MethodPut, "/lists/:id/items/:item_id", testUserID, h.updateItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPut,
				"/lists/10/items/"+tt.param,
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if called != tt.wantCall {
				t.Fatalf("service called = %v, want %v", called, tt.wantCall)
			}
			if !called {
				return
			}

			if gotUserID != testUserID {
				t.Errorf("service got userID %d, want %d", gotUserID, testUserID)
			}
			if tt.param == "42" && gotItem.ID != 42 {
				t.Errorf("service got item ID %d, want 42", gotItem.ID)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp dto.ItemResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.Done != tt.wantDone {
				t.Errorf("got Done %v, want %v", resp.Done, tt.wantDone)
			}
		})
	}
}

func TestHandler_deleteItem(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		svcErr     error
		wantStatus int
		wantCall   bool
	}{
		{
			name:       "success",
			param:      "42",
			wantStatus: http.StatusNoContent,
			wantCall:   true,
		},
		{
			name:       "invalid item id",
			param:      "abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			param:      "999",
			svcErr:     domain.ErrItemNotFound,
			wantStatus: http.StatusNotFound,
			wantCall:   true,
		},
		{
			name:       "service error",
			param:      "42",
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUserID, gotItemID int

			svc := &mockItemService{
				t: t,
				deleteFn: func(ctx context.Context, uID, iID int) error {
					called = true
					gotUserID, gotItemID = uID, iID
					return tt.svcErr
				},
			}

			h := &Handler{item: svc}
			r := setupRouter(t, http.MethodDelete, "/lists/:id/items/:item_id", testUserID, h.deleteItem)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/lists/10/items/"+tt.param, nil))

			if w.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if called != tt.wantCall {
				t.Fatalf("service called = %v, want %v", called, tt.wantCall)
			}
			if !called {
				return
			}

			if gotUserID != testUserID {
				t.Errorf("service got userID %d, want %d", gotUserID, testUserID)
			}
			if tt.param == "42" && gotItemID != 42 {
				t.Errorf("service got itemID %d, want 42", gotItemID)
			}

			if tt.wantStatus == http.StatusNoContent && w.Body.Len() != 0 {
				t.Errorf("expected empty body, got %q", w.Body.String())
			}
		})
	}
}
