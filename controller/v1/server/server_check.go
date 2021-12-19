package server

import "github.com/gin-gonic/gin"

func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(200)
	}
}
