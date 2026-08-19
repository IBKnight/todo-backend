package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", userTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, domain.ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func (r *AuthRepo) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE username = $1", userTable)

	err := r.db.QueryRow(query, username).
		Scan(&user.Id, &user.Username, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by username: %w", err)
	}

	return user, nil
}
