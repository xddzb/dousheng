package service

import (
	"github.com/xddzb/dousheng/model"
	"github.com/xddzb/dousheng/utils"
)

// PostVideo 投稿视频
func PostVideo(userId int64, videoName, coverName, title string) error {
	return NewPostVideoFlow(userId, videoName, coverName, title).Do()
}

func NewPostVideoFlow(userId int64, videoName, coverName, title string) *PostVideoFlow {
	return &PostVideoFlow{
		videoName: videoName,
		coverName: coverName,
		userId:    userId,
		title:     title,
	}
}

type PostVideoFlow struct {
	videoName string
	coverName string
	title     string
	userId    int64

	video *model.Video
}

func (f *PostVideoFlow) Do() error {
	f.prepareParam()

	if err := f.publish(); err != nil {
		return err
	}
	return nil
}

// 准备好参数
func (f *PostVideoFlow) prepareParam() {
	f.videoName = utils.GetFileUrl(f.videoName)
	f.coverName = utils.GetFileUrl(f.coverName)
}

// 组合并添加到数据库
func (f *PostVideoFlow) publish() error {
	video := &model.Video{
		UserInfoId: f.userId,
		PlayUrl:    f.videoName,
		CoverUrl:   f.coverName,
		Title:      f.title,
	}
	return model.NewVideoDAO().AddVideo(video)
}
