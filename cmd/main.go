package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
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
	r.GET("/pai/home", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.tmpl", nil)
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
