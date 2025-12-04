package handlers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/restful-app/pkg/web"
	"github.com/gin-gonic/gin"
)

// ShowLandingPage renders the HTML landing page
func (h *Handler) ShowLandingPage(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", web.LandingPage)
}

// ShowHealthPage renders the HTML health status page
func (h *Handler) ShowHealthPage(c *gin.Context) {
	// Check if the client accepts HTML
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", web.HealthPage)
}
