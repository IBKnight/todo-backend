package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/dto"
)

func TestHandler_signUp(t *testing.T) {
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
			body:       `{"name":"Ivan","username":"ivan","password":"secret123"}`,
			svcID:      42,
			wantStatus: http.StatusCreated,
			wantCall:   true,
		},
		{
			name:       "missing name",
			body:       `{"username":"ivan","password":"secret123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing username",
			body:       `{"name":"Ivan","password":"secret123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       `{"name":"Ivan","username":"ivan"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed json",
			body:       `{"username":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "username already taken",
			body:       `{"name":"Ivan","username":"ivan","password":"secret123"}`,
			svcErr:     domain.ErrUserExists,
			wantStatus: http.StatusConflict,
			wantCall:   true,
		},
		{
			name:       "service error",
			body:       `{"name":"Ivan","username":"ivan","password":"secret123"}`,
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUser domain.User

			svc := &mockAuthService{
				t: t,
				createUserFn: func(user domain.User) (int, error) {
					called = true
					gotUser = user
					if tt.svcErr != nil {
						return 0, tt.svcErr
					}
					return tt.svcID, nil
				},
			}

			h := &Handler{auth: svc}
			r := setupRouter(t, http.MethodPost, "/auth/sign-up", 0, h.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if called != tt.wantCall {
				t.Fatalf("service called = %v, want %v", called, tt.wantCall)
			}

			if strings.Contains(w.Body.String(), "secret123") {
				t.Errorf("password leaked into response: %s", w.Body.String())
			}

			if !called {
				return
			}

			if gotUser.Name != "Ivan" || gotUser.Username != "ivan" {
				t.Errorf("service got %+v", gotUser)
			}
			if gotUser.Password != "secret123" {
				t.Errorf("service got password %q, want raw password", gotUser.Password)
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			var resp dto.UserResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.ID != tt.svcID {
				t.Errorf("got ID %d, want %d", resp.ID, tt.svcID)
			}
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcToken   string
		svcErr     error
		wantStatus int
		wantCall   bool
	}{
		{
			name:       "success",
			body:       `{"username":"ivan","password":"secret123"}`,
			svcToken:   "eyJhbGciOiJIUzI1NiJ9.test.signature",
			wantStatus: http.StatusOK,
			wantCall:   true,
		},
		{
			name:       "missing username",
			body:       `{"password":"secret123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       `{"username":"ivan"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed json",
			body:       `{"username":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong password",
			body:       `{"username":"ivan","password":"wrongpass"}`,
			svcErr:     domain.ErrInvalidCredentials,
			wantStatus: http.StatusUnauthorized,
			wantCall:   true,
		},
		{
			name:       "user does not exist",
			body:       `{"username":"ivan","password":"secret123"}`,
			svcErr:     domain.ErrUserNotFound,
			wantStatus: http.StatusUnauthorized,
			wantCall:   true,
		},
		{
			name:       "service error",
			body:       `{"username":"ivan","password":"secret123"}`,
			svcErr:     errInternal,
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotUsername, gotPassword string

			svc := &mockAuthService{
				t: t,
				generateTokenFn: func(username, password string) (string, error) {
					called = true
					gotUsername, gotPassword = username, password
					if tt.svcErr != nil {
						return "", tt.svcErr
					}
					return tt.svcToken, nil
				},
			}

			h := &Handler{auth: svc}
			r := setupRouter(t, http.MethodPost, "/auth/sign-in", 0, h.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if called != tt.wantCall {
				t.Fatalf("service called = %v, want %v", called, tt.wantCall)
			}

			body := w.Body.String()
			if strings.Contains(body, "secret123") || strings.Contains(body, "wrongpass") {
				t.Errorf("password leaked into response: %s", body)
			}

			if !called {
				return
			}

			if gotUsername != "ivan" {
				t.Errorf("service got username %q, want %q", gotUsername, "ivan")
			}
			if gotPassword == "" {
				t.Error("service got empty password")
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp dto.SignInResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if resp.Token != tt.svcToken {
				t.Errorf("got token %q, want %q", resp.Token, tt.svcToken)
			}
		})
	}
}

// Отказ по разным причинам должен выглядеть для клиента одинаково,
// иначе перебором можно выяснить, какие логины заняты.
func TestHandler_signIn_DoesNotLeakUserExistence(t *testing.T) {
	respond := func(t *testing.T, svcErr error) (int, string) {
		t.Helper()

		svc := &mockAuthService{
			t: t,
			generateTokenFn: func(string, string) (string, error) {
				return "", svcErr
			},
		}

		h := &Handler{auth: svc}
		r := setupRouter(t, http.MethodPost, "/auth/sign-in", 0, h.signIn)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/sign-in",
			strings.NewReader(`{"username":"ivan","password":"anypassword"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		return w.Code, w.Body.String()
	}

	codeNoUser, bodyNoUser := respond(t, domain.ErrUserNotFound)
	codeBadPass, bodyBadPass := respond(t, domain.ErrInvalidCredentials)

	if codeNoUser != codeBadPass {
		t.Errorf("different status: missing user %d, wrong password %d", codeNoUser, codeBadPass)
	}
	if bodyNoUser != bodyBadPass {
		t.Errorf("different body: missing user %q, wrong password %q", bodyNoUser, bodyBadPass)
	}
}
