package model

import (
	"errors"
	"sync"
	"time"
)

type Video struct {
	Id            int64     `json:"id" gorm:"id,omitempty"`
	Author        UserInfo  `json:"author" gorm:"-"`
	PlayUrl       string    `json:"play_url" gorm:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url" gorm:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count" gorm:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count" gorm:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite" gorm:"is_favorite,omitempty"`
	Title         string    `json:"title,omitempty" gorm:"title,omitempty""`
	CreatedTime   time.Time `json:"-" gorm:"created_time,omitempty""`
	UpdatedTime   time.Time `json:"-" gorm:"updated_time,omitempty""`
}

type VideoDAO struct {
}

var (
	videoDAO  *VideoDAO
	videoOnce sync.Once
)

func NewVideoDAO() *VideoDAO {
	videoOnce.Do(func() {
		videoDAO = new(VideoDAO)
	})
	return videoDAO
}

// QueryVideoListByLimitAndTime  返回按投稿时间倒序的视频列表，并限制为最多limit个
func (v *VideoDAO) QueryVideoListByLimitAndTime(limit int, latestTime time.Time, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByLimit videoList 空指针")
	}
	return db.Model(&Video{}).Where("created_time<?", latestTime).
		Order("created_time ASC").Limit(limit).
		Select([]string{"id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title", "created_time", "updated_time"}).
		Find(videoList).Error
}
