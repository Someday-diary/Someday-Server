package controller

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func CreatePost(c *gin.Context) {
	email := c.Request.Header.Get("email")
	key := c.Request.Header.Get("key")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	diaries, err := model.UnmarshalCreatePostRequest(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	db, err := model.ConnectMYSQL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer db.Close()

	tx, _ := db.Begin()
	defer tx.Rollback()

	for _, diary := range diaries.Diaries {
		newUuid, _ := uuid.NewV4()
		postId := newUuid.String()
		_, err = tx.Exec("insert into post values (?, ?, HEX(AES_ENCRYPT(?, ?)), ?)", postId, email, diary.Contents, key, diary.Date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			log.Panic(err)
			return
		}
		for _, tag := range diary.Tags {
			_, err = tx.Exec("insert into tag values (?, HEX(AES_ENCRYPT(?, ?)))", postId, tag.TagName, key)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err.Error(),
				})
				log.Panic(err)

				return
			}
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func GetPostByID(c *gin.Context) {
	c.JSON(200, "call GPBI")

}

func GetAllPost(c *gin.Context) {
	c.JSON(200, "call GAP")

}

func GetPostByTag(c *gin.Context) {
	c.JSON(200, "call GPBT")

}

func EditPost(c *gin.Context) {
	c.JSON(200, "call EP")
}

func DeletePost(c *gin.Context) {
	c.JSON(200, "call DP")

}
