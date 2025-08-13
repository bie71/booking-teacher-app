package middleware

import (
	"auth/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(jwtConfig *service.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn().Msg("No authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn().Msg("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwtConfig.ValidateToken(tokenString)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Warn().Msg("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Warn().Msg("User ID not found in token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("isAdmin", claims["role"] == "admin")
		c.Set("token", tokenString)
		c.Next()
	}
}
