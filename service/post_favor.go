package service

import (
	"errors"
	"github.com/xddzb/dousheng/model"
	"github.com/xddzb/dousheng/redis"
)

const (
	ADD      = 1
	SUBTRACT = 2
)

func PostFavorState(userId, videoId, action int64) error {
	return NewPostFavorStateFlow(userId, videoId, action).Do()
}

type PostFavorStateFlow struct {
	userId  int64
	videoId int64
	action  int64
}

func NewPostFavorStateFlow(userId, videoId, action int64) *PostFavorStateFlow {
	return &PostFavorStateFlow{
		userId:  userId,
		videoId: videoId,
		action:  action,
	}
}

func (f *PostFavorStateFlow) Do() error {
	//检查输入参数
	var err error
	if err = f.checkParameter(); err != nil {
		return err
	}
	//执行操作
	switch f.action {
	case ADD:
		err = f.AddOperation()
	case SUBTRACT:
		err = f.SubtractOperation()
	default:
		return errors.New("未定义的操作")
	}
	return err
}

func (f *PostFavorStateFlow) checkParameter() error {
	if !model.NewUserInfoDAO().IsUserExistById(f.userId) {
		return errors.New("用户不存在")
	}
	if f.action != ADD && f.action != SUBTRACT {
		return errors.New("未定义的行为")
	}
	return nil
}

// 点赞操作 更新m数据库和redis缓存
func (f *PostFavorStateFlow) AddOperation() error {
	//视频点赞数目+1
	err := model.NewVideoDAO().AddOneFavorByUserIdAndVideoId(f.userId, f.videoId)
	if err != nil {
		return errors.New("不要重复点赞")
	}
	//对应的用户是否点赞更新到内存中
	redis.NewProxyIndexMap().UpdateVideoFavorState(f.userId, f.videoId, true)
	return nil
}

// 取消点赞操作
func (f *PostFavorStateFlow) SubtractOperation() error {
	//视频点赞数目-1
	err := model.NewVideoDAO().SubOneFavorByUserIdAndVideoId(f.userId, f.videoId)
	if err != nil {
		return errors.New("点赞数目已经为0")
	}
	//对应的用户是否点赞更新到内存中
	redis.NewProxyIndexMap().UpdateVideoFavorState(f.userId, f.videoId, false)
	return nil
}
