package middleware

import (
	"context"
	"gopher-gateway-service/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager jwt.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		accessToken := parts[1]
		claims, err := jwtManager.ValidateAccessToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 1. Store in Gin context (for Gin handlers)
		c.Set("public_user_id", claims.PublicUserId)
		
		// 2. Store in Request context (for httputil.ReverseProxy)
		ctx := context.WithValue(c.Request.Context(), "public_user_id", claims.PublicUserId)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
