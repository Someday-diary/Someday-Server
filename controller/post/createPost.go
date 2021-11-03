package post

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type CreatePostRequest struct {
	Tags []struct {
		TagName string `json:"tag"`
	} `json:"tags"`
	Contents string `json:"contents"`
	Date     string `json:"date"`
	ID       string `json:"id"`
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

		email := c.Request.Header.Get("email")
		fmt.Println(email)

		aes := lib.CreateCipher(key)

		t, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 409,
			})
			return
		}

		e, err := aes.Encrypt(req.Contents)
		fmt.Println(e)
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
