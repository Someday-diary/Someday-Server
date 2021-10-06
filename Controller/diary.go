package Controller

import "github.com/gin-gonic/gin"

func CreatePost(c *gin.Context) {
	c.JSON(200, "call CP")
}

func GetPostByID(c *gin.Context) {
	c.JSON(200, "call GPBI")

}

func GetAllPost(c *gin.Context) {
	c.JSON(200, "call GAP")

}

func GetPostByTag(c *gin.Context) {
	c.JSON(200, "call GPBT")

}

func EditPost(c *gin.Context) {
	c.JSON(200, "call EP")
}

func DeletePost(c *gin.Context) {
	c.JSON(200, "call DP")

}
