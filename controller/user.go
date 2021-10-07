package controller

import "github.com/gin-gonic/gin"

func SignUp(c *gin.Context) {
	c.JSON(200, "call SI")
}

func Login(c *gin.Context) {
	c.JSON(200, "call LI")
}
