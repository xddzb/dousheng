package user

import (
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/service"
	"net/http"
)

type RegisterResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	rawVal, _ := c.Get("password")
	password, ok := rawVal.(string)
	if !ok {
		c.JSON(http.StatusOK, RegisterResponse{
			StatusCode: 1,
			StatusMsg:  "密码解析出错",
		})
		return
	}

	registerResponse, err := service.AddUser(username, password)

	if err != nil {
		c.JSON(http.StatusOK, RegisterResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     registerResponse.UserId,
		Token:      registerResponse.Token,
	})
}
