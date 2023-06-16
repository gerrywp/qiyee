package main

import (
	"errors"
	"fmt"
	"net/http"

	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/configs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	db := configs.GetDB()
	var user models.User
	result := db.Where("user_name=?", "ii").First(&user)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("登录成功")
	} else {
		fmt.Println("登录失败")
	}
	r.StaticFS("/static", http.Dir("../web/static"))
	r.LoadHTMLGlob("../web/views/**/*")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pai",
		})
	})
	r.GET("/pai/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.tmpl", nil)
	})
	r.POST("/pai/login", func(ctx *gin.Context) {

	})
	r.GET("/pai/home", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home2.tmpl", nil)
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
