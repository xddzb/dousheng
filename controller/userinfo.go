package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/model"
	"log"
	"strconv"
	"strings"
)

func GetUserInfo(c *gin.Context) {
	log.Print("进入GetUserInfo函数")
	userId := c.Param("id")
	userId = strings.TrimLeft(userId, ":,")
	println(userId)
	userinfoDAO := model.NewUserInfoDAO()
	var userInfo model.UserInfo
	useridInt, _ := strconv.ParseInt(userId, 10, 64)

	err := userinfoDAO.QueryUserInfoById(useridInt, &userInfo)
	if err != nil {
		//return err
	}
	c.JSON(200, userInfo)
}
