package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Video struct {
	Id            int64       `json:"id" gorm:"id,omitempty"`
	UserInfoId    int64       `json:"-" gorm:"user_info_id,omitempty"`
	Author        UserInfo    `json:"author" gorm:"-"`
	PlayUrl       string      `json:"play_url" gorm:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url" gorm:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count" gorm:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count" gorm:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite" gorm:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty" gorm:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	CreatedTime   time.Time   `json:"-" gorm:"created_time,omitempty"`
	UpdatedTime   time.Time   `json:"-" gorm:"updated_time,omitempty"`
}
type UserFavorVideos struct {
	UserInfoId int64
	VideoId    int64
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

// AddVideo 添加视频
func (v *VideoDAO) AddVideo(video *Video) error {
	if video == nil {
		return errors.New("AddVideo video 空指针")
	}
	return db.Create(video).Error
}

func (v *VideoDAO) QueryVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList 空指针")
	}
	return db.Where("user_info_id=?", userId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(videoList).Error
}

func (v *VideoDAO) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId count 空指针")
	}
	return db.Model(&Video{}).Where("user_info_id=?", userId).Count(count).Error
}

func (v *VideoDAO) IsVideoExistById(id int64) bool {
	var video Video
	if err := db.Where("id=?", id).Select("id").First(&video).Error; err != nil {
		log.Println(err)
	}
	if video.Id == 0 {
		return false
	}
	return true
}

// AddOneFavorByUserIdAndVideoId 增加一个赞
func (v *VideoDAO) AddOneFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		//gorm.Expr("quantity - ?", 1)
		if err := db.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			log.Println(err)
			return err
		}
		userfavorvideo := &UserFavorVideos{
			UserInfoId: userId,
			VideoId:    videoId,
		}
		//插入操作
		if err := db.Create(&userfavorvideo).Error; err != nil {
			return err
		}
		return nil
	})
}

// SubOneFavorByUserIdAndVideoId 减少一个赞
func (v *VideoDAO) SubOneFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		//db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
		//执行减一操作前先检查favorite_count是否大于0
		if err := db.Model(&Video{}).Where("id = ? AND favorite_count >= ?", videoId, 0).
			Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			log.Println(err)
			return err
		}
		//删除操作
		//db.Where("name = ?", "jinzhu").Delete(&email)
		if err := db.Where("UserInfoId = ? AND VideoId = ?", userId, videoId).
			Delete(&UserFavorVideos{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (v *VideoDAO) QueryFavorVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryFavorVideoListByUserId videoList 空指针")
	}
	//多表查询，左连接得到结果，再映射到数据
	if err := db.Raw("SELECT v.* FROM user_favor_videos u , videos v WHERE u.user_info_id = ? AND u.video_id = v.id", userId).Scan(videoList).Error; err != nil {
		return err
	}
	//如果id为0，则说明没有查到数据
	if len(*videoList) == 0 || (*videoList)[0].Id == 0 {
		return errors.New("点赞列表为空")
	}
	return nil
}
