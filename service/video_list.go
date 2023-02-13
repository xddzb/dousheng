package service

import (
	"errors"
	"github.com/xddzb/dousheng/model"
)

type List struct {
	Videos []*model.Video `json:"video_list,omitempty"`
}

func QueryVideoListByUserId(userId int64) (*List, error) {
	return NewQueryVideoListByUserIdFlow(userId).Do()
}

func NewQueryVideoListByUserIdFlow(userId int64) *QueryVideoListByUserIdFlow {
	return &QueryVideoListByUserIdFlow{userId: userId}
}

type QueryVideoListByUserIdFlow struct {
	userId int64

	videos []*model.Video

	videoList *List
}

func (f *QueryVideoListByUserIdFlow) Do() (*List, error) {
	if err := f.checkNum(); err != nil {
		return nil, err
	}
	if err := f.packData(); err != nil {
		return nil, err
	}
	return f.videoList, nil
}

func (f *QueryVideoListByUserIdFlow) checkNum() error {
	//检查userI对应的用户信息是否存在
	if !model.NewUserInfoDAO().IsUserExistById(f.userId) {
		return errors.New("用户不存在")
	}

	return nil
}

// 注意：Video由于在数据库中没有存储作者信息，所以需要手动填充
func (f *QueryVideoListByUserIdFlow) packData() error {
	//获取该用户投稿的视频列表
	err := model.NewVideoDAO().QueryVideoListByUserId(f.userId, &f.videos)
	if err != nil {
		return err
	}
	//作者信息查询
	var userInfo model.UserInfo
	err = model.NewUserInfoDAO().QueryUserInfoById(f.userId, &userInfo)
	if err != nil {
		return err
	}
	//填充信息(Author和IsFavorite字段
	//p := redis.NewProxyIndexMap()
	for i := range f.videos {
		f.videos[i].Author = userInfo
		f.videos[i].IsFavorite = true
		//f.videos[i].IsFavorite = p.GetVideoFavorState(f.userId, f.videos[i].Id)
	}

	f.videoList = &List{Videos: f.videos}

	return nil
}
