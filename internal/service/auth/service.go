package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type AuthService struct {
	repo     service.AuthorizationRepository
	secret   []byte
	tokenTTL time.Duration
}

func NewService(repo service.AuthorizationRepository, secret []byte, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:     repo,
		secret:   secret,
		tokenTTL: tokenTTL,
	}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}

	user.Password = string(hash)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username string, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)

	if err != nil {
		logrus.Info(user, err, password, username)
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Info(user, err, password, username)
		return "", domain.ErrInvalidCredentials
	}

	now := time.Now()

	claims := tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   strconv.Itoa(user.Id),
		},
		UserID: user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *AuthService) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok || !token.Valid {
		return 0, domain.ErrInvalidToken
	}

	return claims.UserID, nil

}
