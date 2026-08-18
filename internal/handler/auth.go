package handler

import (
	"net/http"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signIn(ctx *gin.Context) {
	var req SignInRequest

	if err := ctx.BindJSON(&req); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.auth.GenerateToken(req.Username, req.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, SignInResponse{
		Token: token,
	})

}

func (h *Handler) signUp(ctx *gin.Context) {
	var req SignUpRequest

	if err := ctx.BindJSON(&req); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user := domain.User{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	}

	id, err := h.auth.CreateUser(user)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, UserResponse{
		ID: id,
	})

}
