package main

import (
	"net/http"
	"path/filepath"

	"gerry.wang/qiyee/api/router"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.StaticFS("/static", http.Dir("../web/static"))
	r.HTMLRender = loadTemplates("../web/views/**")
	//r.LoadHTMLGlob("../web/views/**/*")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pai",
		})
	})
	store := cookie.NewStore([]byte("paigo$"))
	r.Use(sessions.Sessions("qiyee", store))
	router.SetupRouter(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// 非模板嵌套
	htmls, err := filepath.Glob(templatesDir + "/*.tmpl")
	if err != nil {
		panic(err.Error())
	}
	for _, html := range htmls {
		r.AddFromGlob(filepath.Base(html), html)
	}
	//布局模板
	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	//layout.tmpl必须在切片第一个位置
	for i, j := 0, len(layouts)-1; i < j; i, j = i+1, j-1 {
		layouts[i], layouts[j] = layouts[j], layouts[i]
	}
	if err != nil {
		panic(err.Error())
	}
	//嵌套的内容模板
	includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
	if err != nil {
		panic(err.Error())
	}
	for _, include := range includes {
		// 文件名称
		baseName := filepath.Base(include)
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(baseName, files...)
	}
	return r
}
