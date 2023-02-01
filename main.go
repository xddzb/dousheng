package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/controller/test"
	"github.com/xddzb/dousheng/controller/user"
	"github.com/xddzb/dousheng/controller/video"
	"github.com/xddzb/dousheng/model"
)

func main() {
	//初始化数据库
	model.InitDb()
	//初始化引擎配置
	r := gin.Default()
	// public文件存放静态资源
	r.Static("/static", "./public")
	//构建路由
	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", video.Feed)
	apiRouter.GET("/user/", user.GetUserInfo)
	apiRouter.POST("/user/register/", user.Register)
	apiRouter.POST("/user/login/", user.Login)
	apiRouter.POST("/publish/action/", video.Publish)
	apiRouter.GET("/publish/list/", video.PublishList)

	r.GET("/test", test.Test)
	r.GET("/getuserinfo:id", user.GetUserInfo)
	//启动服务
	r.Run(":8087")
}
