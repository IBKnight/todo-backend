package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(ctx *gin.Context) {
	id, _ := ctx.Get(userCtx)
	ctx.JSON(http.StatusOK, map[string]any{
		"id": id,
	})

}
func (h *Handler) getAllLists(ctx *gin.Context) {

}

func (h *Handler) getListById(ctx *gin.Context) {

}

func (h *Handler) updateList(ctx *gin.Context) {

}

func (h *Handler) deleteList(ctx *gin.Context) {

}
