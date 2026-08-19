package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter(t *testing.T, method, path string, userID int, h gin.HandlerFunc) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Handle(method, path, func(c *gin.Context) {
		if userID > 0 {
			c.Set(userCtx, userID)
		}
		h(c)
	})
	return r
}
