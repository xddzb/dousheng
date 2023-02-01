package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/controller"
	"github.com/xddzb/dousheng/model"
)

func main() {
	//初始化数据库
	model.InitDb()
	//初始化引擎配置
	r := gin.Default()
	//构建路由
	r.GET("/test", controller.Test)
	r.GET("/getuserinfo:id", controller.GetUserInfo)
	r.GET("/getvideo:id", controller.Feed)
	//启动服务
	r.Run()
}
