package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartethnet/rustun-dashboard/internal/model"
)

// BasicAuth creates a Basic Authentication middleware
func BasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()

		if !ok {
			unauthorized(c)
			return
		}

		// Use constant time comparison to prevent timing attacks
		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(username)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(password)) == 1

		if !userMatch || !passMatch {
			unauthorized(c)
			return
		}

		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", `Basic realm="Authorization Required"`)
	c.JSON(http.StatusUnauthorized, model.ErrorResponseWithCode(
		http.StatusUnauthorized,
		"Unauthorized",
		"Invalid credentials",
	))
	c.Abort()
}
