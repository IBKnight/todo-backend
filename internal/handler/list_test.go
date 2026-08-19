package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/dto"
	"golang.org/x/net/context"
)

const testUserID = 7

var errInternal = errors.New("db is down")

func TestHandler_createList(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcID      int
		svcErr     error
		wantStatus int
		wantCall   bool
	}{
		{
			name:       "success",
			body:       `{"title":"work","description":"tasks"}`,
			svcID:      42,
			wantStatus: http.StatusOK,
			wantCall:   true,
		},
		{
			name:       "empty title",
			body:       `{"title":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed json",
			body:       `{"title":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			body:       `{"title":"work"}`,
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUserID int
			var gotList domain.TodoList

			svc := &mockListService{
				t: t,
				createFn: func(ctx context.Context, userId int, list domain.TodoList) (int, error) {
					called = true
					gotUserID, gotList = userId, list
					if tt.svcErr != nil {
						return 0, tt.svcErr
					}
					return tt.svcID, nil
				},
			}

			h := &Handler{list: svc}
			r := setupRouter(t, http.MethodPost, "/lists", testUserID, h.createList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/lists", strings.NewReader(tt.body))
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
			if tt.svcErr == nil && gotList.Title != "work" {
				t.Errorf("service got title %q, want %q", gotList.Title, "work")
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp dto.CreatedListResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.ID != tt.svcID {
				t.Errorf("got ID %d, want %d", resp.ID, tt.svcID)
			}
		})
	}
}

func TestHandler_getAllLists(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		want := []domain.TodoList{
			{ID: 1, Title: "work"},
			{ID: 2, Title: "home"},
		}

		var gotUserID int
		svc := &mockListService{
			t: t,
			getAllFn: func(ctx context.Context, userId int) ([]domain.TodoList, error) {
				gotUserID = userId
				return want, nil
			},
		}

		h := &Handler{list: svc}
		r := setupRouter(t, http.MethodGet, "/lists", testUserID, h.getAllLists)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("got status %d, want 200, body: %s", w.Code, w.Body.String())
		}
		if gotUserID != testUserID {
			t.Errorf("service got userID %d, want %d", gotUserID, testUserID)
		}

		var resp dto.ListsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(resp.Data) != 2 {
			t.Fatalf("got %d lists, want 2", len(resp.Data))
		}
		if resp.Data[0].Title != "work" {
			t.Errorf("got title %q, want %q", resp.Data[0].Title, "work")
		}
	})

	t.Run("empty list serializes as array", func(t *testing.T) {
		svc := &mockListService{
			t: t,
			getAllFn: func(context.Context, int) ([]domain.TodoList, error) {
				return nil, nil
			},
		}

		h := &Handler{list: svc}
		r := setupRouter(t, http.MethodGet, "/lists", testUserID, h.getAllLists)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists", nil))

		if strings.Contains(w.Body.String(), "null") {
			t.Errorf("expected empty array, got %s", w.Body.String())
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockListService{
			t: t,
			getAllFn: func(context.Context, int) ([]domain.TodoList, error) {
				return nil, errInternal
			},
		}

		h := &Handler{list: svc}
		r := setupRouter(t, http.MethodGet, "/lists", testUserID, h.getAllLists)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists", nil))

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("got status %d, want 500", w.Code)
		}
	})

	t.Run("no user id in context", func(t *testing.T) {
		svc := &mockListService{t: t} // getAllFn не задан — вызов уронит тест

		h := &Handler{list: svc}
		r := setupRouter(t, http.MethodGet, "/lists", 0, h.getAllLists)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists", nil))

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("got status %d, want 500", w.Code)
		}
	})
}

func TestHandler_getListById(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		svcList    domain.TodoList
		svcErr     error
		wantStatus int
		wantCall   bool
	}{
		{
			name:       "success",
			param:      "42",
			svcList:    domain.TodoList{ID: 42, Title: "work", Description: "tasks"},
			wantStatus: http.StatusOK,
			wantCall:   true,
		},
		{
			name:       "not found",
			param:      "999",
			svcErr:     domain.ErrListNotFound,
			wantStatus: http.StatusNotFound,
			wantCall:   true,
		},
		{
			name:       "invalid id",
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
			var gotListID, gotUserID int

			svc := &mockListService{
				t: t,
				getByIDFn: func(ctx context.Context, listId, userId int) (domain.TodoList, error) {
					called = true
					gotListID, gotUserID = listId, userId
					if tt.svcErr != nil {
						return domain.TodoList{}, tt.svcErr
					}
					return tt.svcList, nil
				},
			}

			h := &Handler{list: svc}
			r := setupRouter(t, http.MethodGet, "/lists/:id", testUserID, h.getListById)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/lists/"+tt.param, nil))

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
			if tt.param == "42" && gotListID != 42 {
				t.Errorf("service got listID %d, want 42", gotListID)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp dto.ListResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.ID != tt.svcList.ID || resp.Title != tt.svcList.Title {
				t.Errorf("got %+v, want ID=%d Title=%q", resp, tt.svcList.ID, tt.svcList.Title)
			}
		})
	}
}

func TestHandler_updateList(t *testing.T) {
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
			param:      "42",
			body:       `{"title":"updated","description":"new"}`,
			wantStatus: http.StatusOK,
			wantCall:   true,
		},
		{
			name:       "invalid id",
			param:      "abc",
			body:       `{"title":"updated"}`,
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
			body:       `{"title":"updated"}`,
			svcErr:     domain.ErrListNotFound,
			wantStatus: http.StatusNotFound,
			wantCall:   true,
		},
		{
			name:       "service error",
			param:      "42",
			body:       `{"title":"updated"}`,
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUserID int
			var gotList domain.TodoList

			svc := &mockListService{
				t: t,
				updateFn: func(ctx context.Context, userId int, list domain.TodoList) (domain.TodoList, error) {
					called = true
					gotUserID, gotList = userId, list
					if tt.svcErr != nil {
						return domain.TodoList{}, tt.svcErr
					}
					return list, nil
				},
			}

			h := &Handler{list: svc}
			r := setupRouter(t, http.MethodPut, "/lists/:id", testUserID, h.updateList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/lists/"+tt.param, strings.NewReader(tt.body))
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
			if tt.param == "42" && gotList.ID != 42 {
				t.Errorf("service got list ID %d, want 42", gotList.ID)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp dto.ListResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.Title != "updated" {
				t.Errorf("got title %q, want %q", resp.Title, "updated")
			}
		})
	}
}

func TestHandler_deleteList(t *testing.T) {
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
			name:       "invalid id",
			param:      "abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			param:      "999",
			svcErr:     domain.ErrListNotFound,
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
			var gotUserID, gotListID int

			svc := &mockListService{
				t: t,
				removeFn: func(ctx context.Context, userId, listId int) error {
					called = true
					gotUserID, gotListID = userId, listId
					return tt.svcErr
				},
			}

			h := &Handler{list: svc}
			r := setupRouter(t, http.MethodDelete, "/lists/:id", testUserID, h.deleteList)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/lists/"+tt.param, nil))

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
			if tt.param == "42" && gotListID != 42 {
				t.Errorf("service got listID %d, want 42", gotListID)
			}

			if tt.wantStatus == http.StatusNoContent && w.Body.Len() != 0 {
				t.Errorf("expected empty body, got %q", w.Body.String())
			}
		})
	}
}
