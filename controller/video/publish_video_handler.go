package video

import (
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/service"
	"github.com/xddzb/dousheng/utils"
	"net/http"
	"path/filepath"
)

type PublishResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

var (
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
	pictureIndexMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)

// 发布视频并截取第一帧画面作为封面
func Publish(c *gin.Context) {
	//准备参数
	rawId, _ := c.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		PublishVideoError(c, "解析UserId出错")
		return
	}
	title := c.PostForm("title")
	form, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}
	//不支持多文件上传
	file := form.File["data"][0]

	suffix := filepath.Ext(file.Filename)    //得到后缀
	if _, ok := videoIndexMap[suffix]; !ok { //判断是否为视频格式
		PublishVideoError(c, "不支持的视频格式")
		return
	}
	name := utils.NewFileName(userId) //根据userId得到唯一的文件名
	videoFilename := name + suffix
	videoSavePath := filepath.Join("./public", videoFilename)
	err = c.SaveUploadedFile(file, videoSavePath)
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}
	pictureSavePath := filepath.Join("./public", name)
	//截取一帧画面作为封面
	_, err = utils.GetSnapshot(videoSavePath, pictureSavePath, 1)
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}
	//数据库持久化
	erro := service.PostVideo(userId, videoFilename, name+".jpg", title)
	if erro != nil {
		PublishVideoError(c, erro.Error())
		return
	}
	PublishVideoOk(c, file.Filename+"上传成功")

}

func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, PublishResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func PublishVideoOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, PublishResponse{
		StatusCode: 0,
		StatusMsg:  msg,
	})
}
