package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func Feed(c *gin.Context) {
	userId := c.Param("id")
	userId = strings.TrimLeft(userId, ":,")
	println(userId)
	data := "video1"
	c.JSON(200, data)
}
