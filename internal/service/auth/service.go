package auth

import "github.com/IBKnight/todo-backend/internal/repository"

type Authorization struct {
	repo *repository.AuthRepo
}

func NewService(repo *repository.AuthRepo) *Authorization {
	return &Authorization{
		repo: repo,
	}
}
