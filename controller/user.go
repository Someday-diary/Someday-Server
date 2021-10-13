package controller

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	loginRequest, err := model.UnmarshalSignRequest(body)
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
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	var verity string

	err = db.Get(&verity, "select if ((select verify from user where email = ?) IS NULL or (select verify from user where email = ?) = 'N', 'N', 'Y');", loginRequest.Email, loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}

	if verity == "N" {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "이메일 인증 안함 응애",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Pwd), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	_, err = db.Exec("update user set pwd = ?, agree = ? where email = ?", string(hash), loginRequest.Agree, loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	b := make([]byte, 32)
	_, _ = rand.Read(b)
	key := fmt.Sprintf("%x", b)[:32]

	client, err := model.ConnectRedis(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		_, _ = db.Exec("delete FROM user where email = ?", loginRequest.Email)
		return
	}

	_, err = client.Set(context.Background(), loginRequest.Email, key, 0).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		_, _ = db.Exec("delete FROM user where email = ?", loginRequest.Email)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func Login(c *gin.Context) {

}

func SendVerityEmail(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	verityRequest, err := model.UnmarshalEmailVerityRequest(body)
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
	_, err = db.Exec("insert into user (email) value (?)", verityRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	code := string(b)

	client, err := model.ConnectRedis(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	_, err = client.Set(context.Background(), verityRequest.Email, code, time.Minute*10).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// TODO: 메일 주소 파면 바꾸기 and 디자인 적용하기
	templateData := map[string]string{
		"Code": code,
	}

	r := model.NewRequest(verityRequest.Email, os.Getenv("smtp_id"), "[오늘하루] 회원가입을 위한 인증번호를 알려드립니다", "")
	err = r.ParseTemplate("email.html", templateData)

	err = r.SendEmail()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
}

func Verity(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	req, err := model.UnmarshalEmailVerityConfirmRequest(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	d, err := model.ConnectMYSQL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	r, err := model.ConnectRedis(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	code, err := r.Get(context.Background(), req.Email).Result()
	if err != nil {
		code = ""
	}

	if code != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "코드 안맞아 ~",
		})
		return
	}

	_, err = d.Exec("update user set verify = 'Y' where email = ?", req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	_, err = r.Del(context.Background(), req.Email).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		_, _ = d.Exec("update user set verify = 'N' where email = ?", req.Email)
		return
	}

}
