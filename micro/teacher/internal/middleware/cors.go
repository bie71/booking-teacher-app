package middleware

import (
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a CORS middleware with appropriate settings
func CORSMiddleware() gin.HandlerFunc {
	// Baca env ALLOWED_ORIGINS dan pisah pakai koma
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	// Biar rapi, trim spasi tiap origin
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Cek apakah origin termasuk dalam list
		if slices.Contains(allowedOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept, Accept-Encoding, Accept-Language, Connection, Host")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
