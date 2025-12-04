package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	msgAuthHeaderRequired    = "Authorization header is required"
	msgAuthHeaderFormat      = "Authorization header format must be Bearer {token}"
	msgInvalidToken          = "Invalid token"
	msgInvalidTokenClaims    = "Invalid token claims"
	msgInvalidUserIDInClaims = "Invalid user ID in token claims"
	msgTokenExpired          = "Token has expired"
	msgUserNotFound          = "User not found"
	msgInternalError         = "An internal error occurred"
)

func (app *application) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgAuthHeaderRequired})
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgAuthHeaderFormat})
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// require HMAC signing to prevent algorithm confusion
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(app.jwtSecret), nil
		})

		if err != nil {
			var e *jwt.ValidationError
			if errors.As(err, &e) && e.Errors&jwt.ValidationErrorExpired != 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgTokenExpired})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgInvalidToken})
			}
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgInvalidToken})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgInvalidTokenClaims})
			return
		}

		// support numeric (float64) and string user_id claim types
		var uid int
		switch v := claims["user_id"].(type) {
		case float64:
			uid = int(v)
		case string:
			i, convErr := strconv.Atoi(v)
			if convErr != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgInvalidUserIDInClaims})
				return
			}
			uid = i
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgInvalidUserIDInClaims})
			return
		}

		user, err := app.models.Users.Get(uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msgInternalError})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msgUserNotFound})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
