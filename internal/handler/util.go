package handler

import "github.com/gin-gonic/gin"

func getUserID(c *gin.Context) (int, bool) {
	id, exists := c.Get(userCtx)
	if !exists {
		return 0, false
	}
	idInt, ok := id.(int)
	return idInt, ok
}
