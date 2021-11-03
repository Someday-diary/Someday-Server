package post

import (
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
)

type EditPostRequest struct {
	ID       string `json:"id"`
	Contents string `json:"contents"`
	Date     string `json:"date"`
	Tags     []struct {
		TagName string `json:"tag"`
	} `json:"tags"`
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

		for _, diary := range req.Diaries {
			var post model.Post
			model.DB.Model(&post).Where("id = ?", diary.Id).Updates(&model.Post{
				Contents: diary.Contents,
			})

			model.DB.Where("post_id = ?", diary.Id).Delete(&model.Tag{})
			for _, tag := range diary.Tags {
				t := model.Tag{
					PostID:  diary.Id,
					TagName: tag.TagName,
				}
				model.DB.Create(&t)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
