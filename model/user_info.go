package model

import (
	"errors"
	"sync"
)

var (
	ErrIvdPtr = errors.New("空指针错误")
)

type UserInfo struct {
	Id            int64  `json:"id" gorm:"id,omitempty"`
	Name          string `json:"name" gorm:"name,omitempty"`
	FollowCount   int64  `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow" gorm:"is_follow,omitempty"`
}

type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once
)

func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}
func (u *UserInfoDAO) QueryUserInfoById(userId int64, userinfo *UserInfo) error {
	if userinfo == nil {
		return ErrIvdPtr
	}
	//DB.Where("id=?",userId).First(userinfo)
	db.Where("id=?", userId).Select([]string{"id", "name", "follow_count", "follower_count",
		"is_follow"}).First(userinfo)
	//id为零值，说明sql执行失败
	if userinfo.Id == 0 {
		return errors.New("该用户不存在")
	}
	return nil
}
