package controller

import "github.com/gin-gonic/gin"

// 测试web服务是否启动成功
func Test(c *gin.Context) {
	c.JSON(200, gin.H{"code": 200, "msg": "i see you"})
}
