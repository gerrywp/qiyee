package main

import (
	"net/http"

	"gerry.wang/qiyee/api/router"
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
	router.SetupRouter(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
