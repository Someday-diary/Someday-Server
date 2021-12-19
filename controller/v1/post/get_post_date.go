package post

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPostByDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		type response struct {
			Code int   `json:"code,omitempty"`
			Post *post `json:"post,omitempty"`
		}
		res := new(response)
		res.Post = new(post)
		res.Post.Tags = new([]tag)

		year, ok := c.GetQuery("year")
		if ok == false {
			res.Code = 400
			res.Post = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}
		month, ok := c.GetQuery("month")
		if ok == false {
			res.Code = 400
			res.Post = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}

		day, ok := c.GetQuery("day")
		if ok == false {
			res.Code = 400
			res.Post = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}

		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		cipher := lib.CreateCipher(key)

		date := fmt.Sprintf("%s-%s-%s", year, month, day)
		var post model.Post
		err := model.DB.Where("email = ? and date_format(created_at, '%Y %m %d') = "+
			"date_format(?, '%Y %m %d')", email, date).First(&post).Error
		if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
			panic(err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Code = 110
			res.Post = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}

		model.DB.Find(&post.Tag, "post_id = ?", post.ID)

		res.Post.PostID = post.ID
		res.Post.Date = post.CreatedAt.Format("2006-01-02")
		res.Post.Contents, err = cipher.Decrypt(post.Contents)
		if err != nil {
			panic(err)
		}

		for _, t := range post.Tag {
			temp := tag{}
			temp.TagName, err = cipher.Decrypt(t.TagName)
			if err != nil {
				panic(err)
			}
			*res.Post.Tags = append(*res.Post.Tags, temp)
		}

		res.Code = 200
		c.JSON(http.StatusOK, res)
	}
}
