package middleware

import (
	"booking/internal/config"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// simple in-memory blacklist (token -> expiry unix)
var tokenBlacklist = struct {
	m map[string]int64
	sync.RWMutex
}{m: make(map[string]int64)}

// BlacklistToken adds token until expiry ts
func BlacklistToken(token string, exp int64) {
	tokenBlacklist.Lock()
	defer tokenBlacklist.Unlock()
	tokenBlacklist.m[token] = exp
}

// IsBlacklisted checks if token is in blacklist and not expired
func IsBlacklisted(token string) bool {
	tokenBlacklist.RLock()
	exp, ok := tokenBlacklist.m[token]
	tokenBlacklist.RUnlock()
	if !ok {
		return false
	}
	if time.Now().Unix() > exp {
		// cleanup
		tokenBlacklist.Lock()
		delete(tokenBlacklist.m, token)
		tokenBlacklist.Unlock()
		return false
	}
	return true
}

func AuthMiddleware(jwtConfig *config.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Simpan ke context
		c.Set("user_id", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))
		c.Set("token", tokenStr)

		c.Next()
	}
}
