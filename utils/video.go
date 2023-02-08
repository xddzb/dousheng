package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/xddzb/dousheng/config"
	"github.com/xddzb/dousheng/model"
	"log"
	"os"
	"strings"
	"time"
)

func GetFileUrl(fileName string) string {
	base := fmt.Sprintf("http://%s:%d/static/%s", config.Info.Server.IP, config.Info.Server.Port, fileName)
	return base
}

// NewFileName 根据userId+用户发布的视频数量连接成独一无二的文件名
func NewFileName(userId int64) string {
	var count int64
	err := model.NewVideoDAO().QueryVideoCountByUserId(userId, &count)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d", userId, count)
}

// FillVideoFields 填充每个视频的作者信息 因为作者与视频的一对多关系，数据库中存下的是作者的id
// 当userId>0时，我们判断当前为登录状态，其余情况为未登录状态，则不需要填充IsFavorite字段
func FillVideoFields(userId int64, video *model.Video) (*time.Time, error) {
	if video == nil {
		return nil, errors.New("util.FillVideoFields video为空")
	}
	dao := model.NewUserInfoDAO()
	latestTime := video.CreatedTime //获取最近的投稿时间
	//添加作者信息，以及is_follow状态
	var userInfo model.UserInfo
	err := dao.QueryUserInfoById(video.UserInfoId, &userInfo)
	if err != nil {
		return nil, err
	}
	video.Author = userInfo

	return &latestTime, nil
}

/*
videoPath: 视频文件地址
snapshotPath: 生成图片的地址
frameNum: 获取第几帧
*/
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".jpg")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	names := strings.Split(snapshotPath, "\\")
	snapshotName = names[len(names)-1] + ".jpg"
	return
}
