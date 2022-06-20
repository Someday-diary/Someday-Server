package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model/dao"
	"github.com/Someday-diary/Someday-Server/model/database"
	"github.com/Someday-diary/Someday-Server/model/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(dto.CreatePost)
		db := database.ConnectDB()
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
		p := dao.Post{
			ID:        req.ID,
			Email:     email,
			Contents:  e,
			CreatedAt: t,
			Tag:       nil,
		}
		db.Create(&p)

		for _, tag := range req.Tags {
			e, err := aes.Encrypt(tag.TagName)
			if err != nil {
				panic(err)
			}
			t := dao.Tag{
				PostID:  req.ID,
				TagName: e,
			}
			db.Create(&t)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"post_id": req.ID,
		})
	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.ConnectDB()

		var post dao.Post

		postID := c.Param("post_id")
		email := c.GetHeader("email")
		err := db.First(&post, "id = ?", postID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 110,
			})
			return
		} else if post.Email != email {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 111,
			})
			return
		}

		err = db.Delete(&dao.Post{ID: postID}).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}

func EditPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.ConnectDB()
		req := new(dto.EditPost)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		postID := c.Param("post_id")
		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		cipher := lib.CreateCipher(key)

		var post dao.Post
		err = db.First(&post, "id = ?", postID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 110,
			})
			return
		} else if post.Email != email {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 111,
			})
			return
		}

		e, err := cipher.Encrypt(req.Contents)
		if err != nil {
			panic(err)
		}

		db.Model(&post).Where("id = ?", postID).Updates(&dao.Post{
			Contents: e,
		})

		db.Where("post_id = ?", postID).Delete(&dao.Tag{})
		for _, tag := range req.Tags {
			e, err := cipher.Encrypt(tag.TagName)
			if err != nil {
				panic(err)
			}

			t := dao.Tag{
				PostID:  postID,
				TagName: e,
			}
			db.Create(&t)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}

func GetPostByDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.ConnectDB()
		type response struct {
			Code int       `json:"code,omitempty"`
			Post *dto.Post `json:"post,omitempty"`
		}
		res := new(response)
		res.Post = new(dto.Post)
		res.Post.Tags = new([]dto.Tag)

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
		var post dao.Post
		err := db.Where("email = ? and date_format(created_at, '%Y %m %d') = "+
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

		db.Find(&post.Tag, "post_id = ?", post.ID)

		res.Post.PostID = post.ID
		res.Post.Date = post.CreatedAt.Format("2006-01-02")
		res.Post.Contents, err = cipher.Decrypt(post.Contents)
		if err != nil {
			panic(err)
		}

		for _, t := range post.Tag {
			temp := dto.Tag{}
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

func GetPostByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.ConnectDB()

		type response struct {
			Code int       `json:"code,omitempty"`
			Post *dto.Post `json:"post,omitempty"`
		}
		res := new(response)
		res.Post = new(dto.Post)
		res.Post.Tags = new([]dto.Tag)

		id := c.Param("post_id")
		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		cipher := lib.CreateCipher(key)

		var post dao.Post
		var n int64
		err := db.First(&post, "id = ?", id).Count(&n).Error
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
		db.Find(&post.Tag, "post_id = ?", id)

		res.Post.PostID = post.ID
		res.Post.Date = post.CreatedAt.Format("2006-01-02")
		res.Post.Contents, err = cipher.Decrypt(post.Contents)
		if err != nil {
			panic(err)
		}

		for _, t := range post.Tag {
			temp := dto.Tag{}
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

func GetPostByMonth() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.ConnectDB()

		type response struct {
			Code  int         `json:"code,omitempty"`
			Posts *[]dto.Post `json:"posts,omitempty"`
		}
		res := new(response)
		res.Posts = new([]dto.Post)

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
		var posts []dao.Post
		db.Select("post.*").Where("email = ? and date_format(created_at, '%Y %m') = "+
			"date_format(?, '%Y %m')", email, date).Find(&posts)

		if len(posts) == 0 {
			res.Code = 110
			res.Posts = nil
			c.JSON(http.StatusBadRequest, res)
			return
		}

		for _, p := range posts {
			temp := dto.Post{
				PostID: p.ID,
				Date:   p.CreatedAt.Format("2006-01-02"),
			}
			*res.Posts = append(*res.Posts, temp)
		}

		res.Code = 200
		c.JSON(http.StatusOK, res)
	}
}

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.ConnectDB()

		type response struct {
			Code  int         `json:"code,omitempty"`
			Posts *[]dto.Post `json:"posts,omitempty"`
		}
		var err error
		res := new(response)
		res.Posts = new([]dto.Post)

		req := c.QueryArray("tags")

		email := c.GetHeader("email")
		key := c.GetHeader("secret_key")

		cipher := lib.CreateCipher(key)

		var posts []dao.Post
		var n int64

		if len(req) != 0 {
			for i, tag := range req {
				req[i], err = cipher.Encrypt(tag)
				if err != nil {
					panic(err)
				}
			}

			err = db.Raw("SELECT post.* FROM tag JOIN post ON post.id = tag.post_id WHERE post.email = ? and tag.tag_name in "+
				"(?) GROUP BY tag.post_id having (count(tag.tag_name) = ?)",
				email, req, len(req)).First(&posts).Order("post.created_at desc").Count(&n).Error
		} else {
			err = db.Raw("select * from post where email = ?", email).First(&posts).Count(&n).Error
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
			temp := dto.Post{
				PostID: p.ID,
				Date:   p.CreatedAt.Format("2006-01-02"),
				Tags:   new([]dto.Tag),
			}

			err = db.Find(&posts[i].Tag, "post_id = ?", p.ID).Error
			if err != nil {
				panic(err)
			}

			temp.Contents, err = cipher.Decrypt(p.Contents)
			if err != nil {
				panic(err)
			}

			for _, t := range posts[i].Tag {
				tempTag := dto.Tag{}
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
