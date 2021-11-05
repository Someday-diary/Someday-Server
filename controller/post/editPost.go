package post

import (
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
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

		key := c.GetHeader("secret_key")

		aes := lib.CreateCipher(key)

		var post model.Post
		e, err := aes.Encrypt(req.Contents)
		if err != nil {
			panic(err)
		}

		model.DB.Model(&post).Where("id = ?", req.ID).Updates(&model.Post{
			Contents: e,
		})

		model.DB.Where("post_id = ?", req.ID).Delete(&model.Tag{})
		for _, tag := range req.Tags {
			e, err := aes.Encrypt(tag.TagName)
			if err != nil {
				panic(err)
			}

			t := model.Tag{
				PostID:  req.ID,
				TagName: e,
			}
			model.DB.Create(&t)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
