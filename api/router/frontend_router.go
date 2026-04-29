package router

import (
	"net/http"

	"gerry.wang/qiyee/api/service"
	"github.com/gin-gonic/gin"
)

// SetupFrontendRouter 配置前端网站路由
func SetupFrontendRouter(r *gin.Engine) *gin.Engine {
	// 配置首页路由
	r.GET("/", index)
	r.GET("/about", about)
	return r
}

// 首页处理函数
func index(ctx *gin.Context) {
	var as = service.NewAbout()
	var ps = service.NewProder()
	var bs = service.NewBanner()
	var ss = service.NewSite()
	var ns = service.NewNews()

	// 获取数据
	about := as.GetAbout()
	products := ps.GetProds()
	banners := bs.GetBanners()
	site := ss.GetSite()
	news, _ := ns.GetNewsByPage(1, 5)

	// 准备模板数据
	data := gin.H{
		"About":    *about,
		"Products": products,
		"Banners":  banners,
		"News":     news,
		"Site":     site,
	}

	ctx.HTML(http.StatusOK, "index.tmpl", data)
}

func about(ctx *gin.Context) {
	var as = service.NewAbout()
	about := as.GetAbout()
	ctx.HTML(http.StatusOK, "about.tmpl", *about)
}
