package post

import (
	"net/http"
	"time"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Tags []struct {
		TagName string `json:"tag" binding:"required"`
	} `json:"tags"`
	Contents string `json:"contents" binding:"required"`
	Date     string `json:"date" binding:"required"`
	ID       string `json:"id" binding:"required"`
}

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(CreatePostRequest)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		keys := make(map[string]struct {
			TagName string `json:"tag" binding:"required"`
		})
		tags := make([]struct {
			TagName string `json:"tag" binding:"required"`
		}, 0)
		for _, v := range req.Tags {
			if _, ok := keys[v.TagName]; ok {
				continue
			} else {
				keys[v.TagName] = v
				tags = append(tags, v)
			}
		}

		req.Tags = tags

		aes := lib.CreateCipher(key)

		t, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 409,
			})
			return
		}

		e, err := aes.Encrypt(req.Contents)
		if err != nil {
			panic(err)
		}
		p := model.Post{
			ID:        req.ID,
			Email:     email,
			Contents:  e,
			CreatedAt: t,
			Tag:       nil,
		}
		model.DB.Create(&p)

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
			"code":    200,
			"post_id": req.ID,
		})
	}
}
