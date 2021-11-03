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

type PRes struct {
	PostID   string `json:"post_id"`
	Contents string `json:"contents"`
	Email    string `json:"email"`
	Date     string `json:"date"`

	Tag []Tag `json:"tag"`
}

type Tag struct {
	TagName string `json:"tag_name"`
}

func GetPostByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("post_id")
		email := c.GetHeader("email")

		var secretKey model.Secret
		err := model.DB.Find(&secretKey, "email = ?", email).Error
		if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
			panic(err)
		}

		key, err := lib.SystemCipher.Decrypt(secretKey.SecretKey)
		if err != nil {
			panic(err)
		}

		aes := lib.CreateCipher(key)

		var post model.Post
		var n int64
		err = model.DB.First(&post, "id = ?", id).Count(&n).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}

		if n == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 110,
			})
			return
		}
		if post.Email != c.GetHeader("email") {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 111,
			})
			return
		}
		model.DB.Find(&post.Tag, "post_id = ?", id)

		res := PRes{}

		res.PostID = post.ID
		res.Date = post.CreatedAt.Format("2006-01-02")
		res.Email = post.Email
		res.Contents, err = aes.Decrypt(post.Contents)
		if err != nil {
			panic(err)
		}

		for _, tag := range post.Tag {
			temp := Tag{}
			temp.TagName, err = aes.Decrypt(tag.TagName)
			if err != nil {
				panic(err)
			}
			res.Tag = append(res.Tag, temp)
		}

		c.JSON(http.StatusOK, res)
	}
}

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.QueryArray("tags")
		var err error

		email := c.Request.Header.Get("email")
		key := c.GetHeader("secret_key")

		aes := lib.CreateCipher(key)

		var posts []model.Post

		if len(req) != 0 {
			for i, tag := range req {
				req[i], err = aes.Encrypt(tag)
				if err != nil {
					panic(err)
				}
			}

			err = model.DB.Raw("SELECT post.* FROM tag JOIN post ON post.id = tag.post_id WHERE post.email = ? and tag.tag_name in "+
				"(?) GROUP BY tag.post_id having (count(tag.tag_name) = ?)",
				email, req, len(req)).First(&posts).Error
		} else {
			err = model.DB.Raw("select * from post where email = ?", email).First(&posts).Error
		}

		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				panic(err)
			}
		}

		var res []PRes

		for i, post := range posts {
			temp := PRes{
				PostID: post.ID,
				Email:  post.Email,
				Date:   post.CreatedAt.Format("2006-01-02"),
			}

			err = model.DB.Find(&posts[i].Tag, "post_id = ?", post.ID).Error
			if err != nil {
				panic(err)
			}

			temp.Contents, err = aes.Decrypt(post.Contents)
			if err != nil {
				panic(err)
			}

			for _, tag := range posts[i].Tag {
				tempTag := Tag{}
				tempTag.TagName, err = aes.Decrypt(tag.TagName)
				if err != nil {
					panic(err)
				}
				temp.Tag = append(temp.Tag, tempTag)
			}

			res = append(res, temp)
		}
		c.JSON(200, res)
	}
}

type postResponse struct {
	Date   string `json:"date"`
	PostID string `json:"post_id"`
}

func GetPostByDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, ok := c.GetQuery("year")
		if ok == false {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
		}
		month, ok := c.GetQuery("month")
		if ok == false {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
		}

		email := c.GetHeader("email")

		date := fmt.Sprintf("%s-%s-1", year, month)
		var posts []model.Post
		model.DB.Select("id, created_at").Where("email = ? and date_format(created_at, '%Y %m') = "+
			"date_format(?, '%Y %m')", email, date).Find(&posts)

		if len(posts) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 110,
			})
			return
		}

		var res []postResponse

		fmt.Println(posts)
		for _, post := range posts {
			temp := postResponse{}
			t := post.CreatedAt.Format("2006-01-02")
			temp.Date = t
			temp.PostID = post.ID

			res = append(res, temp)
		}

		c.JSON(http.StatusOK, res)
	}
}
