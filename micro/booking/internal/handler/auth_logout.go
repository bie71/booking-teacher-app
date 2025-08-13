package handler

import (
	"net/http"
	"time"
	"booking/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Logout(c *gin.Context) {
	v, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error":"No token in context"})
		return
	}
	exp := time.Now().Add(24*time.Hour).Unix()
	if ev, ok := c.Get("exp"); ok {
		if e64, ok2 := ev.(int64); ok2 { exp = e64 }
	}
	middleware.BlacklistToken(v.(string), exp)
	c.JSON(http.StatusOK, gin.H{"message":"Logged out"})
}
