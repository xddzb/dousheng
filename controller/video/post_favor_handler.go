package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/service"
	"net/http"
	"strconv"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
type ProxyPostFavorHandler struct {
	*gin.Context

	userId  int64
	videoId int64
	action  int64 //1-点赞，2-取消点赞
}

func PostFavorHandler(c *gin.Context) {
	NewProxyPostFavorHandler(c).Do()
}

func NewProxyPostFavorHandler(c *gin.Context) *ProxyPostFavorHandler {
	return &ProxyPostFavorHandler{Context: c}
}

func (p *ProxyPostFavorHandler) Do() {
	//解析UserId videoId action等内容
	if err := p.ParseParameter(); err != nil {
		p.SendError(err.Error())
		return
	}
	//调用service层方法
	err := service.PostFavorState(p.userId, p.videoId, p.action)
	if err != nil {
		p.SendError(err.Error())
		return
	}
	//成功返回
	p.SendOk()
}

func (p *ProxyPostFavorHandler) ParseParameter() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId转int失败")
	}
	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		return err
	}
	p.userId = userId
	p.videoId = videoId
	p.action = actionType
	return nil
}

func (p *ProxyPostFavorHandler) SendError(msg string) {
	p.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: msg})
}

func (p *ProxyPostFavorHandler) SendOk() {
	p.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
}
