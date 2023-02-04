package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/model"
	"log"
	"net/http"
)

type UserResponse struct {
	StatusCode int             `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	User       *model.UserInfo `json:"user"`
}

func GetUserInfo(c *gin.Context) {
	log.Print("进入GetUserInfo函数")
	//得到上层中间件根据token解析的userId
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusOK, UserResponse{
			StatusCode: 1,
			StatusMsg:  "解析userId出错",
		})
		return
	}
	userinfo, err := DoQueryUserInfoByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		User:       userinfo,
	})
	return

}

func DoQueryUserInfoByUserId(rawId interface{}) (*model.UserInfo, error) {
	userId, err := rawId.(int64)
	if !err {
		return nil, errors.New("解析userId失败")
	}
	//由于得到userinfo不需要组装model层的数据，所以直接调用model层的接口
	userinfoDAO := model.NewUserInfoDAO()
	var userInfo model.UserInfo
	ok := userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if ok != nil {
		return nil, ok
	}

	return &userInfo, nil
}

