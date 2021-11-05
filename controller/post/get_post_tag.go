package post

import (
	"errors"
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		type response struct {
			Code  int     `json:"code,omitempty"`
			Posts *[]post `json:"posts,omitempty"`
		}
		res := new(response)
		res.Posts = new([]post)

		req := c.QueryArray("tags")
		var err error

		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		cipher := lib.CreateCipher(key)

		var posts []model.Post
		var n int64

		if len(req) != 0 {
			for i, tag := range req {
				req[i], err = cipher.Encrypt(tag)
				if err != nil {
					panic(err)
				}
			}

			err = model.DB.Raw("SELECT post.* FROM tag JOIN post ON post.id = tag.post_id WHERE post.email = ? and tag.tag_name in "+
				"(?) GROUP BY tag.post_id having (count(tag.tag_name) = ?)",
				email, req, len(req)).First(&posts).Count(&n).Error
		} else {
			err = model.DB.Raw("select * from post where email = ?", email).First(&posts).Count(&n).Error
		}

		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				panic(err)
			}
		}

		if n == 0 {
			res.Code = 110
			res.Posts = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}

		for i, p := range posts {
			temp := post{
				PostID: p.ID,
				Date:   p.CreatedAt.Format("2006-01-02"),
				Tags:   new([]tag),
			}

			err = model.DB.Find(&posts[i].Tag, "post_id = ?", p.ID).Error
			if err != nil {
				panic(err)
			}

			temp.Contents, err = cipher.Decrypt(p.Contents)
			if err != nil {
				panic(err)
			}

			for _, t := range posts[i].Tag {
				tempTag := tag{}
				tempTag.TagName, err = cipher.Decrypt(t.TagName)
				if err != nil {
					panic(err)
				}
				*temp.Tags = append(*temp.Tags, tempTag)
			}

			*res.Posts = append(*res.Posts, temp)
		}
		res.Code = 200
		c.JSON(200, res)
	}
}
