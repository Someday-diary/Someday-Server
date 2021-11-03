package post

import (
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
)

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var post model.Post
		err := model.DB.Delete(&post, "id = ?", c.Param("post_id")).Error
		if err != nil {
			panic(err)
		}

		if n == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 110,
			})
			return
		}

		if email := c.GetHeader("email"); email != post.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 111,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
