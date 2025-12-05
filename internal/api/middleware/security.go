package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds common security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Protects against XSS attacks
		c.Header("X-XSS-Protection", "1; mode=block")
		// Prevents the browser from interpreting files as a different MIME type
		c.Header("X-Content-Type-Options", "nosniff")
		// Prevents the site from being embedded in an iframe (clickjacking protection)
		c.Header("X-Frame-Options", "DENY")
		// Controls how much referrer information is included with requests
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Enforces HTTPS (HSTS) - 1 year
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Content Security Policy (Basic)
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")

		c.Next()
	}
}
