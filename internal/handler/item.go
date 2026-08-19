package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/IBKnight/todo-backend/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id")
		return
	}

	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	created, err := h.item.Create(c.Request.Context(), userID, listID, req.ToDomain())
	if err != nil {
		h.handleError(c, err, "create item", userID)
		return
	}

	c.Header("Location", fmt.Sprintf("/api/items/%d", created.ID))
	c.JSON(http.StatusCreated, dto.NewItemResponse(created))
}

func (h *Handler) getAllItems(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id")
		return
	}

	items, err := h.item.GetAll(c.Request.Context(), userID, listID)
	if err != nil {
		h.handleError(c, err, "get all items", userID)
		return
	}

	c.JSON(http.StatusOK, dto.NewItemsResponse(items))
}

func (h *Handler) getItemById(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.item.GetByID(c.Request.Context(), userID, itemID)
	if err != nil {
		h.handleError(c, err, "get item by id", userID)
		return
	}

	c.JSON(http.StatusOK, dto.NewItemResponse(item))
}

func (h *Handler) updateItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id")
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	updated, err := h.item.Update(c.Request.Context(), userID, req.ToDomain(itemID))
	if err != nil {
		h.handleError(c, err, "update item", userID)
		return
	}

	c.JSON(http.StatusOK, dto.NewItemResponse(updated))
}

func (h *Handler) deleteItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.item.Delete(c.Request.Context(), userID, itemID); err != nil {
		h.handleError(c, err, "delete item", userID)
		return
	}

	c.Status(http.StatusNoContent)
}
