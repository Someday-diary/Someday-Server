package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	}
}
