package service

import (
	"github.com/xddzb/dousheng/model"
	"strconv"
	"sync"
	"time"
)

// MaxVideoNum 每次最多返回的视频流数量
const (
	MaxVideoNum = 30
)

type VideoInfo struct {
	Videos   []*model.Video `json:"video_list"`
	NextTime int64          `json:"next_time"`
}
type QueryVideoInfoFlow struct {
	//controller层调用时传的参数 需要对参数进行检查
	userId    int64
	Timestamp string
	//中间变量 用于接收model层返回的数据
	videos      []*model.Video
	nextTime    int64
	lastestTime time.Time
	//存放经过处理的信息 返回给controller层
	videoInfo *VideoInfo
}

func QueryFeedVideoList(userId int64, Timestamp string) (*VideoInfo, error) {
	return newQueryVideoInfoFlow(userId, Timestamp).Do()
}

func newQueryVideoInfoFlow(userId int64, Timestamp string) *QueryVideoInfoFlow {
	return &QueryVideoInfoFlow{
		Timestamp: Timestamp,
		userId:    userId,
	}
}

func (f *QueryVideoInfoFlow) Do() (*VideoInfo, error) {
	//参数校验
	if err := f.checkParam(); err != nil {
		return nil, err
	}
	//准备数据
	if err := f.prepareInfo(); err != nil {
		return nil, err
	}
	//组装实体
	if err := f.packPageInfo(); err != nil {
		return nil, err
	}
	return f.videoInfo, nil
}

func (f *QueryVideoInfoFlow) checkParam() error {
	//将前端传回的字符串转为time.Time类型
	intTime, err := strconv.ParseInt(f.Timestamp, 10, 64)
	if err != nil {
		return err
	}
	var latestTime time.Time
	latestTime = time.Unix(intTime, 0)
	//若lastestTime为零 将其设置为本地当前时间
	if latestTime.IsZero() {
		latestTime = time.Now()
	}
	f.lastestTime = latestTime
	//上层通过把userId置零，表示userId不存在或不需要
	if f.userId > 0 {
		//这里说明userId是有效的，可以定制性的做一些登录用户的专属视频推荐
	}
	return nil
}

func (f *QueryVideoInfoFlow) prepareInfo() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		model.NewVideoDAO().QueryVideoListByLimitAndTime(MaxVideoNum, f.lastestTime, &f.videos)
	}()
	wg.Wait() //等待信息从model层返回
	//设置nexttime
	f.nextTime = time.Now().Unix()
	return nil
}

func (f *QueryVideoInfoFlow) packPageInfo() error {
	f.videoInfo = &VideoInfo{
		Videos:   f.videos,
		NextTime: f.nextTime,
	}
	return nil
}
