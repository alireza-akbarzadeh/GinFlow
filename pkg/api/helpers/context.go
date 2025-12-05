package helpers

import (
	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"github.com/gin-gonic/gin"
)

// GetUserFromContext retrieves the authenticated user from gin context
func GetUserFromContext(c *gin.Context) *models.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return nil
	}
	user, ok := contextUser.(*models.User)
	if !ok {
		return nil
	}
	return user
}

// SetUserInContext sets the authenticated user in gin context
func SetUserInContext(c *gin.Context, user *models.User) {
	c.Set("user", user)
}
