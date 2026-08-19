package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTaskNotFound       = errors.New("task not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrListNotFound       = errors.New("list not found")
	ErrItemNotFound       = errors.New("item not found")
	ErrValidation         = errors.New("validation failed")
	ErrUserExists         = errors.New("username already taken")
)
