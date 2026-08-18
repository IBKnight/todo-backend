package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const userCtx = "UserId"

func (h *Handler) UserIdentity(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")

	if header == "" {
		newErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userID, err := h.auth.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid token")
		return

	}

	ctx.Set(userCtx, userID)

}
