package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/config"
	"github.com/xddzb/dousheng/controller/comment"
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
	apiRouter.GET("/user/", middleware.JWTMidWare(), middleware.UserIdVerify(), user.GetUserInfo)
	apiRouter.POST("/user/register/", middleware.SHAMiddleWare(), user.Register) //对用户密码加密存储
	apiRouter.POST("/user/login/", middleware.SHAMiddleWare(), user.Login)       //对用户密码加密存储
	apiRouter.POST("/publish/action/", middleware.JWTMidWare(), video.Publish)
	apiRouter.GET("/publish/list/", middleware.JWTMidWare(), middleware.UserIdVerify(), video.PublishList)

	//extend 1
	apiRouter.POST("/favorite/action/", middleware.JWTMidWare(), video.PostFavorHandler)
	apiRouter.GET("/favorite/list/", middleware.JWTMidWare(), video.FavorVideoListHandler)
	apiRouter.POST("/comment/action/", middleware.JWTMidWare(), comment.PostCommentHandler)
	apiRouter.GET("/comment/list/", middleware.JWTMidWare(), comment.QueryCommentListHandler)

	////extend 2
	apiRouter.POST("/relation/action/", middleware.JWTMidWare(), user.PostFollowActionHandler)
	apiRouter.GET("/relation/follow/list/", middleware.UserIdVerify(), user.QueryFollowListHandler)
	apiRouter.GET("/relation/follower/list/", middleware.UserIdVerify(), user.QueryFollowerHandler)

	//启动服务
	r.Run(fmt.Sprintf(":%d", config.Info.Port))
}
