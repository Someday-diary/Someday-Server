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

func GetPostByMonth() gin.HandlerFunc {
	return func(c *gin.Context) {
		type response struct {
			Code  int     `json:"code,omitempty"`
			Posts *[]post `json:"posts,omitempty"`
		}
		res := new(response)
		res.Posts = new([]post)

		year, ok := c.GetQuery("year")
		if ok == false {
			res.Code = 400
			res.Posts = nil
			c.JSON(http.StatusBadRequest, res)
		}
		month, ok := c.GetQuery("month")
		if ok == false {
			res.Code = 400
			res.Posts = nil
			c.JSON(http.StatusBadRequest, res)
		}

		email := c.GetHeader("email")

		date := fmt.Sprintf("%s-%s-1", year, month)
		var posts []model.Post
		model.DB.Select("post.*").Where("email = ? and date_format(created_at, '%Y %m') = "+
			"date_format(?, '%Y %m')", email, date).Find(&posts)

		if len(posts) == 0 {
			res.Code = 110
			res.Posts = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}

		for _, p := range posts {
			temp := post{
				PostID: p.ID,
				Date:   p.CreatedAt.Format("2006-01-02"),
			}
			*res.Posts = append(*res.Posts, temp)
		}

		res.Code = 200
		c.JSON(http.StatusOK, res)
	}
}

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

		fmt.Println(posts)
		for _, post := range posts {
			temp := postResponse{}
			t := post.CreatedAt.Format("2006-01-02")
			temp.Date = t
			temp.PostID = post.ID

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
