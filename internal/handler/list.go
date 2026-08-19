package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/IBKnight/todo-backend/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) createList(ctx *gin.Context) {
	id, ok := getUserID(ctx)

	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return
	}

	var list dto.CreateListRequest

	if err := ctx.BindJSON(&list); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	listId, err := h.list.CreateList(ctx.Request.Context(), id, list.ToDomain())

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "cannot create list")
		return
	}

	ctx.JSON(http.StatusOK, dto.CreatedListResponse{
		ID: listId,
	})

}
func (h *Handler) getAllLists(ctx *gin.Context) {
	userId, ok := getUserID(ctx)

	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return
	}

	lists, err := h.list.GetUserLists(ctx.Request.Context(), userId)

	if err != nil {
		logrus.Error("get user lists", "err", err, "user_id", userId)
		newErrorResponse(ctx, http.StatusInternalServerError, "failed to get lists")
		return
	}

	ctx.JSON(http.StatusOK, dto.NewListsResponse(lists))

}

func (h *Handler) getListById(ctx *gin.Context) {
	userId, ok := getUserID(ctx)

	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "failed to parse list id")
		return
	}

	list, err := h.list.GetListById(ctx.Request.Context(), listId, userId)

	if err != nil {
		if errors.Is(err, domain.ErrListNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, "list not found")
		}

		logrus.Error("get user lists", "err", err, "user_id", userId)
		newErrorResponse(ctx, http.StatusInternalServerError, "failed to get list")
		return
	}

	ctx.JSON(http.StatusOK, dto.NewListResponse(list))

}

func (h *Handler) updateList(ctx *gin.Context) {
	userId, ok := getUserID(ctx)

	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "failed to parse list id")
		return
	}

	var reqList dto.UpdateListRequest

	if err := ctx.BindJSON(&reqList); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	domList := domain.TodoList{
		ID:          listId,
		Title:       reqList.Title,
		Description: reqList.Description,
	}

	updatedList, err := h.list.UpdateList(ctx.Request.Context(), userId, domList)

	if err != nil {
		if errors.Is(err, domain.ErrListNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, "list not found")
			return
		}

		logrus.Error("get user lists", "err", err, "user_id", userId)
		newErrorResponse(ctx, http.StatusInternalServerError, "failed to get list")
		return
	}

	ctx.JSON(http.StatusOK, dto.ListResponse{
		ID:          updatedList.ID,
		Title:       updatedList.Title,
		Description: updatedList.Description,
	})

}

func (h *Handler) deleteList(ctx *gin.Context) {
	userId, ok := getUserID(ctx)

	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "failed to parse list id")
		return
	}

	if err := h.list.RemoveList(ctx.Request.Context(), userId, listId); err != nil {
		if errors.Is(err, domain.ErrListNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, "list not found")
			return
		}

		logrus.Error("get user lists", "err", err, "user_id", userId)
		newErrorResponse(ctx, http.StatusInternalServerError, "failed to get list")
		return
	}

	ctx.Status(http.StatusNoContent)

}
