package post

import (
	"errors"
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var post model.Post

		postID := c.Param("post_id")
		email := c.GetHeader("email")
		err := model.DB.First(&post, "id = ?", postID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 110,
			})
			return
		} else if post.Email != email {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 111,
			})
			return
		}

		err = model.DB.Delete(&post, "id = ?", c.Param("post_id")).Error
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
