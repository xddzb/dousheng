package video

import (
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/service"
	"net/http"
	"strings"
)

type FormatResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	*service.VideoInfo
}

func Feed(c *gin.Context) {
	rawTimestamp := c.Query("latest_time")
	_, ok := c.GetQuery("token")
	rawTimestamp = strings.TrimLeft(rawTimestamp, ":,")
	//无登录状态
	if !ok {

	}
	videoList, err := service.QueryFeedVideoList(0, rawTimestamp)
	if err != nil {
		c.JSON(http.StatusOK, FormatResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, FormatResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			VideoInfo:  videoList,
		})
	}

}
