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
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "포스트가 없습니다",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	}
}
