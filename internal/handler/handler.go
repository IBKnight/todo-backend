package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth Authorization
	list TodoList
	item TodoItem
}

func NewHandler(
	auth Authorization,
	list TodoList,
	item TodoItem,
) *Handler {
	return &Handler{
		auth: auth,
		list: list,
		item: item,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}

	api := router.Group("/api", h.UserIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("", h.createList)
			lists.GET("", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group("/:id/items")
			{
				items.POST("", h.createItem)
				items.GET("", h.getAllItems)
				items.GET("/:item_id", h.getItemById)
				items.PUT("/:item_id", h.updateItem)
				items.DELETE("/:item_id", h.deleteItem)
			}
		}

	}

	return router
}
