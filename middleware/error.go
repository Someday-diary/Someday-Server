package middleware

import "github.com/gin-gonic/gin"

func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		//	Todo: 에러 핸들링하기
	}
}
