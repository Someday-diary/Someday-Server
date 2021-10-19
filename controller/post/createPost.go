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
	Diaries []struct {
		Tags []struct {
			TagName string `json:"tag"`
		} `json:"tags"`
		Contents string `json:"contents"`
		Date     string `json:"date"`
	} `json:"diaries"`
}

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(CreatePostRequest)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "리퀘 똑바로 날리십쇼..",
			})
			return
		}

		email := c.Request.Header.Get("email")
		fmt.Println(email)

		for _, diary := range req.Diaries {
			id, _ := uuid.NewV4()
			u := id.String()

			t, err := time.Parse("2006-01-02", diary.Date)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":        "시간 똑바로 입력 안하냐?",
					"네놈이 입력한 시간": diary.Date,
				})
				return
			}

			e, err := lib.Cipher.Encrypt(diary.Contents)
			if err != nil {
				_ = c.Error(err)
				return
			}
			p := model.Post{
				ID:        u,
				Email:     email,
				Contents:  e,
				CreatedAt: t,
				Tag:       nil,
			}
			model.DB.Create(&p)

			for _, tag := range diary.Tags {
				e, err := lib.Cipher.Encrypt(tag.TagName)
				if err != nil {
					_ = c.Error(err)
					return
				}
				t := model.Tag{
					PostID:  u,
					TagName: e,
				}
				model.DB.Create(&t)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	}
}
