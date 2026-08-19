package handler

import (
	"errors"
	"net/http"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	user := domain.User{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	}

	id, err := h.auth.CreateUser(user)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserExists):
			newErrorResponse(c, http.StatusConflict, "username already taken")
		case errors.Is(err, domain.ErrValidation):
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			logrus.WithFields(logrus.Fields{
				"op":       "sign up",
				"username": req.Username,
				"error":    err,
			}).Error("failed to create user")
			newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{ID: id})
}

func (h *Handler) signIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	token, err := h.auth.GenerateToken(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
			return
		}

		logrus.WithFields(logrus.Fields{
			"op":       "sign in",
			"username": req.Username,
			"error":    err,
		}).Error("failed to generate token")
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, dto.SignInResponse{Token: token})
}
