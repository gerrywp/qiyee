package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"

	"gerry.wang/qiyee/api/router"
	"gerry.wang/qiyee/middleware"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.StaticFS("/static", http.Dir("./web/static"))
	r.HTMLRender = loadTemplates()
	//r.LoadHTMLGlob("./web/views/**/*")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pai",
		})
	})
	//store := cookie.NewStore([]byte("..1234567.pai-secret-go-key$"))
	store := memstore.NewStore([]byte("..1234567.pai-secret-go-key$"))
	r.Use(sessions.Sessions("qiyee", store))
	middleware.IgnorePaths = []string{
		"/pai/login",
		"/",
		"/about",
		"/products",
		"/news",
		"/news/:id",
		"/contact",
	}
	r.Use(middleware.CheckLogin())
	router.SetupRouter(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 加载后端模板
	loadBackendTemplates(r, "./web/views/backend")

	// 加载前端模板
	loadFrontendTemplates(r, "./web/views/frontend")

	return r
}

// 加载后端模板
func loadBackendTemplates(r multitemplate.Renderer, dir string) {
	// 非模板嵌套
	htmls, err := filepath.Glob(dir + "/*.tmpl")
	if err != nil {
		panic(err.Error())
	}
	for _, html := range htmls {
		r.AddFromGlob(filepath.Base(html), html)
	}

	//布局模板
	layouts, err := filepath.Glob(dir + "/layouts/*.tmpl")
	//layout.tmpl必须在切片第一个位置
	for i, j := 0, len(layouts)-1; i < j; i, j = i+1, j-1 {
		layouts[i], layouts[j] = layouts[j], layouts[i]
	}
	if err != nil {
		panic(err.Error())
	}

	//嵌套的内容模板
	includes, err := filepath.Glob(dir + "/includes/*.tmpl")
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
}

// 加载前端模板
func loadFrontendTemplates(r multitemplate.Renderer, dir string) {
	// 布局模板
	layouts, err := filepath.Glob(dir + "/layouts/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	// 定义安全渲染函数和其他自定义函数
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"stripTags": func(s string) string {
			re := regexp.MustCompile(`(?s)<[^>]*>`)
			return re.ReplaceAllString(s, "")
		},
		"substr": func(s string, start, length int) string {
			if start < 0 || start >= len(s) {
				return ""
			}
			if start+length > len(s) {
				length = len(s) - start
			}
			return s[start : start+length]
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"intRange": func(start, end int) []int {
			var result []int
			for i := start; i < end; i++ {
				result = append(result, i)
			}
			return result
		},
	}

	// 首页模板
	indexTmpl, err := filepath.Glob(dir + "/*.tmpl")
	if err != nil {
		panic(err.Error())
	}
	for _, tmpl := range indexTmpl {
		baseName := filepath.Base(tmpl)
		// 将首页模板与布局模板关联
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, tmpl)
		// 使用 AddFromFilesFuncs 注册自定义函数
		r.AddFromFilesFuncs(baseName, funcMap, files...)
	}
}
