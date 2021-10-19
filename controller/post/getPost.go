package post

import (
	"fmt"
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
)

func GetPostByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("post_id")
		fmt.Println(id)

		var post model.Post

		err := model.DB.First(&post, "id = ?", id).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "일치하는 일기가 없습니다",
			})
			return
		}
		if post.Email != c.GetHeader("email") {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "자신의 일기가 아닙니다.",
			})
			return
		}
		model.DB.Find(&post.Tag, "post_id = ?", id)

		c.JSON(http.StatusOK, post)
	}
}

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		tags := c.QueryArray("tags")

		sql := "select distinct p.id from tag left join post p on p.id = tag.post_id where p.email = ?"

		for i, _ := range tags {
			if i == 0 {
				sql += " and tag_name = ?"
			} else {
				sql += " or tag_name = ?"
			}
		}
		var posts []model.Post
		model.DB.Raw(sql, c.GetHeader("email"), tags[0], tags[1]).Find(&posts)

		c.JSON(200, posts)
	}
}
