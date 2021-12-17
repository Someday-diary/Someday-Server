package post

import (
	"errors"
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tag struct {
	TagName string `json:"tag_name"`
}

type post struct {
	PostID   string `json:"post_id,omitempty"`
	Contents string `json:"contents,omitempty"`
	Date     string `json:"date,omitempty"`

	Tags *[]tag `json:"tags,omitempty"`
}

func GetPostByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		type response struct {
			Code int   `json:"code,omitempty"`
			Post *post `json:"post,omitempty"`
		}
		res := new(response)
		res.Post = new(post)
		res.Post.Tags = new([]tag)

		id := c.Param("post_id")
		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		cipher := lib.CreateCipher(key)

		var post model.Post
		var n int64
		err := model.DB.First(&post, "id = ?", id).Count(&n).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}

		if n == 0 {
			res.Code = 110
			res.Post = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}
		if post.Email != email {
			res.Code = 111
			res.Post = nil
			c.JSON(http.StatusForbidden, res)
			return
		}
		model.DB.Find(&post.Tag, "post_id = ?", id)

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
