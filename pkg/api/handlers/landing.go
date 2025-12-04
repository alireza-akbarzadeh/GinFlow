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
