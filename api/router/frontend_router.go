package router

import (
	"net/http"
	"strconv"

	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/api/service"
	"github.com/gin-gonic/gin"
)

// SetupFrontendRouter 配置前端网站路由
func SetupFrontendRouter(r *gin.Engine) *gin.Engine {
	// 配置首页路由
	r.GET("/", index)
	r.GET("/about", about)
	r.GET("/news", newsList)
	r.GET("/news/:id", newsDetail)
	r.GET("/contact", contact)
	r.GET("/products", productsList)
	r.GET("/products/:id", productDetail)
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
	var ss = service.NewSite()
	about := as.GetAbout()
	site := ss.GetSite()

	// 准备模板数据
	data := gin.H{
		"About": *about,
		"Site":  site,
	}

	ctx.HTML(http.StatusOK, "about.tmpl", data)
}

// 新闻列表页面处理函数
func newsList(ctx *gin.Context) {
	var ns = service.NewNews()
	var ss = service.NewSite()

	// 获取分页参数
	page := ctx.DefaultQuery("page", "1")
	pageNum := 1
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}

	// 获取数据
	news, total := ns.GetNewsByPage(pageNum, 10)
	site := ss.GetSite()

	// 计算分页信息
	pageSize := 10
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	// 准备模板数据
	data := gin.H{
		"News":       news,
		"Site":       site,
		"Page":       pageNum,
		"Total":      total,
		"TotalPages": totalPages,
		"PageSize":   pageSize,
	}

	ctx.HTML(http.StatusOK, "news-list.tmpl", data)
}

// 新闻详情页面处理函数
func newsDetail(ctx *gin.Context) {
	var ns = service.NewNews()
	var ss = service.NewSite()

	// 获取新闻ID
	id := ctx.Param("id")
	newsID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	// 获取新闻详情
	news, err := ns.FindByID(uint(newsID))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	site := ss.GetSite()

	// 准备模板数据
	data := gin.H{
		"News": news,
		"Site": site,
	}

	ctx.HTML(http.StatusOK, "news-detail.tmpl", data)
}

// 联系我们页面处理函数
func contact(ctx *gin.Context) {
	var ss = service.NewSite()
	site := ss.GetSite()

	// 准备模板数据
	data := gin.H{
		"Site": site,
	}

	ctx.HTML(http.StatusOK, "contact.tmpl", data)
}

// 产品列表页面处理函数
func productsList(ctx *gin.Context) {
	var ps = service.NewProder()
	var ss = service.NewSite()

	// 获取分页参数
	page := ctx.DefaultQuery("page", "1")
	pageNum := 1
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}

	// 获取数据
	products, total := ps.GetProdsByPage(pageNum, 12)
	site := ss.GetSite()

	// 计算分页信息
	pageSize := 12
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	// 生成页码列表
	var pageNumbers []int
	startPage := pageNum - 2
	if startPage < 1 {
		startPage = 1
	}
	endPage := pageNum + 2
	if endPage > totalPages {
		endPage = totalPages
	}
	if endPage-startPage < 4 && totalPages > 5 {
		if startPage == 1 {
			endPage = 5
		} else {
			startPage = totalPages - 4
		}
	}

	for i := startPage; i <= endPage; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	// 计算上一页和下一页
	prevPage := pageNum - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := pageNum + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	// 准备模板数据
	data := gin.H{
		"Products":    products,
		"Site":        site,
		"Page":        pageNum,
		"Total":       total,
		"TotalPages":  totalPages,
		"PageSize":    pageSize,
		"PageNumbers": pageNumbers,
		"PrevPage":    prevPage,
		"NextPage":    nextPage,
	}

	ctx.HTML(http.StatusOK, "products.tmpl", data)
}

// 产品详情页面处理函数
func productDetail(ctx *gin.Context) {
	var ps = service.NewProder()
	var ss = service.NewSite()

	// 获取产品ID
	id := ctx.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// 获取产品详情
	product, err := ps.FindByID(uint(productID))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 获取相关产品（排除当前产品）
	allProducts := ps.GetProds()
	var relatedProducts []models.Prod
	for _, p := range allProducts {
		if p.ID != product.ID && len(relatedProducts) < 3 {
			relatedProducts = append(relatedProducts, p)
		}
	}

	site := ss.GetSite()

	// 准备模板数据
	data := gin.H{
		"Product":         product,
		"RelatedProducts": relatedProducts,
		"Site":            site,
	}

	ctx.HTML(http.StatusOK, "products-detail.tmpl", data)
}
