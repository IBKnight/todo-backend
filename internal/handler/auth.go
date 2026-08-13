package handler

import (
	"net/http"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signIn(ctx *gin.Context) {

}

func (h *Handler) signUp(ctx *gin.Context) {
	var input domain.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.auth.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"id": id,
	})

}
