package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/config"
	"github.com/xddzb/dousheng/controller/user"
	"github.com/xddzb/dousheng/controller/video"
	"github.com/xddzb/dousheng/middleware"
	"github.com/xddzb/dousheng/model"
	"log"
)

func main() {
	log.Println("进入main函数")
	//初始化数据库
	model.InitDb()
	log.Println("初始化数据库成功")
	//初始化引擎配置
	r := gin.Default()
	log.Println("初始化引擎配置成功")
	// public文件存放静态资源
	r.Static("/static", "./public")
	log.Println("public文件存放静态资源")
	//构建路由
	apiRouter := r.Group("/douyin")
	log.Println("构建路由成功")
	// basic apis
	apiRouter.GET("/feed/", video.Feed)
	apiRouter.GET("/user/", middleware.JWTMidWare(), user.GetUserInfo)
	apiRouter.POST("/user/register/", middleware.SHAMiddleWare(), user.Register) //对用户密码加密存储
	apiRouter.POST("/user/login/", middleware.SHAMiddleWare(), user.Login)       //对用户密码加密存储
	apiRouter.POST("/publish/action/", middleware.JWTMidWare(), video.Publish)
	apiRouter.GET("/publish/list/", middleware.JWTMidWare(), video.PublishList)

	//测试接口
	//r.GET("/test", test.Test)
	//r.GET("/getuserinfo:id", user.GetUserInfo)
	//启动服务
	r.Run(fmt.Sprintf(":%d", config.Info.Port))

	//_, err := utils.GetSnapshot("./public/test_static.mp4", "./public/test_static", 1)
	//if err != nil {
	//	return
	//}
}
