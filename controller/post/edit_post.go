package post

import (
	"errors"
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EditPostRequest struct {
	Contents string `json:"contents" binding:"required"`
	Tags     []struct {
		TagName string `json:"tag" binding:"required"`
	} `json:"tags" binding:"required"`
}

func EditPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(EditPostRequest)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		postID := c.Param("post_id")
		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		cipher := lib.CreateCipher(key)

		var post model.Post
		err = model.DB.First(&post, "id = ?", postID).Error
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

		e, err := cipher.Encrypt(req.Contents)
		if err != nil {
			panic(err)
		}

		model.DB.Model(&post).Where("id = ?", postID).Updates(&model.Post{
			Contents: e,
		})

		model.DB.Where("post_id = ?", postID).Delete(&model.Tag{})
		for _, tag := range req.Tags {
			e, err := cipher.Encrypt(tag.TagName)
			if err != nil {
				panic(err)
			}

			t := model.Tag{
				PostID:  postID,
				TagName: e,
			}
			model.DB.Create(&t)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
