package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(200)
	}
}

func ReturnError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	}
}
