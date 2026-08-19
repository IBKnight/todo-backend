package handler

import (
	"errors"
	"net/http"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func getUserID(c *gin.Context) (int, bool) {
	id, exists := c.Get(userCtx)
	if !exists {
		return 0, false
	}
	idInt, ok := id.(int)
	return idInt, ok
}

func (h *Handler) handleError(c *gin.Context, err error, op string, userID int) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrListNotFound):
		newErrorResponse(c, http.StatusNotFound, "list not found")
	case errors.Is(err, domain.ErrItemNotFound):
		newErrorResponse(c, http.StatusNotFound, "item not found")
	default:
		logrus.WithFields(logrus.Fields{
			"op":      op,
			"user_id": userID,
			"error":   err,
		}).Error("unhandled error")
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
	}
}
