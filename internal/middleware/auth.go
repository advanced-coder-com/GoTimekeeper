package middleware

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/service"
	"net/http"
	"strings"

	"github.com/advanced-coder-com/go-timekeeper/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserMissingAuthHeader.Error()})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserInvalidAuthHeader.Error()})
			return
		}

		_, claims, err := auth.VerifyJWT(parts[1])
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserTokenInvalid.Error()})
			return
		}

		context.Set("user_id", claims["user_id"])
		context.Next()
	}
}
