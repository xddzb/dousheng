package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserIdVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawId := c.Query("user_id")
		if rawId == "" {
			rawId = c.PostForm("user_id")
		}
		//用户不存在
		if rawId == "" {
			c.JSON(http.StatusOK, Response{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		userId, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
		}
		c.Set("user_id", userId) //若user_id字段数据与token中的ID不同，以此处的user_id为准
		c.Next()
	}
}
