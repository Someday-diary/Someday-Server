package post

import (
	"fmt"
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
)

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
