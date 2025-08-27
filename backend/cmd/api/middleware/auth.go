package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		tkn, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := tkn.Claims.(jwt.MapClaims)
		c.Set("userId", claims["sub"].(string))
		c.Next()
	}
}
