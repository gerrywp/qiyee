package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/api/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PageInfo struct {
	CurrentPage int
	PageSize    int
	TotalCount  int64
	TotalPages  int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
	Pages       []int
}

func buildPageInfo(page, pageSize int, totalCount int64) PageInfo {
	if page < 1 {
		page = 1
	}
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	pages := make([]int, 0, totalPages)
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	return PageInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		Pages:       pages,
	}
}

func login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.tmpl", nil)
}

func doLogin(ctx *gin.Context) {
	var form models.User
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	us := service.NewUser()
	if us.Login(ctx, form.UserName, form.Password, form.IsRemember) {
		//ctx.Redirect(http.StatusFound, "/pai/home")
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "登录成功"})
	} else {
		//ctx.Redirect(http.StatusNotModified, "/pai/login")
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "用户名密码错误!"})
	}
}

func home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	qiyeeData := session.Get("qiyee")
	if qiyeeData != nil {
		var user models.User
		json.Unmarshal([]byte(qiyeeData.(string)), &user)
		ctx.HTML(http.StatusOK, "ihome.tmpl", gin.H{"UserName": user.UserName})
	} else {
		ctx.HTML(http.StatusOK, "ihome.tmpl", nil)
	}
}

func logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("qiyee")
	session.Save()
	ctx.Redirect(http.StatusFound, "/pai/login")
}

func banner(ctx *gin.Context) {
	var bs = service.NewBanner()
	banners := bs.GetBanners()
	m := make(map[string]interface{})
	m["banners"] = banners
	if len(banners) >= 1 {
		m["banner1"] = banners[0]
	}
	if len(banners) >= 2 {
		m["banner2"] = banners[1]
	}
	if len(banners) >= 3 {
		m["banner3"] = banners[2]
	}
	ctx.HTML(http.StatusOK, "banner.tmpl", m)
}

func brand(ctx *gin.Context) {
	var as = service.NewAbout()
	r := as.GetAbout()
	ctx.HTML(http.StatusOK, "brand.tmpl", *r)
}

func brandUpdate(ctx *gin.Context) {
	var content = ctx.PostForm("content")
	var as = service.NewAbout()
	as.Update(content)
	ctx.JSON(http.StatusOK, gin.H{"code": true, "msg": "保存成功"})
}

func bannerUpload(ctx *gin.Context) {
	var bs = service.NewBanner()
	result := bs.Upload(ctx)
	if !result {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "上传失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "上传成功"})
}

// bannerCrop 处理裁切后的图片上传
func bannerCrop(ctx *gin.Context) {
	id := ctx.PostForm("id")
	imageData := ctx.PostForm("imageData")

	if id == "" || imageData == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "参数错误"})
		return
	}

	// TODO: 处理 base64 图片数据并保存
	ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "裁切成功", "thumbUrl": "/static/upload/thumb.jpg"})
}

func prod(ctx *gin.Context) {
	page := 1
	pageParam := ctx.Query("page")
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	const pageSize = 6
	var ps = service.NewProder()
	prods, total := ps.GetProdsByPage(page, pageSize)
	ctx.HTML(http.StatusOK, "prod.tmpl", gin.H{
		"Products": prods,
		"Page":     buildPageInfo(page, pageSize, total),
	})
}

func prodUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("ID"))
	name := ctx.PostForm("Name")
	desc := ctx.PostForm("Desc")
	existingImgUrl := ctx.PostForm("ExistingImgUrl")
	existingThumbUrl := ctx.PostForm("ExistingThumbx150Url")

	prodEntity := models.Prod{
		Name:         name,
		Desc:         desc,
		ImgUrl:       existingImgUrl,
		Thumbx150Url: existingThumbUrl,
	}

	if id > 0 {
		if existing, err := service.NewProder().FindByID(uint(id)); err == nil {
			prodEntity = *existing
			prodEntity.Name = name
			prodEntity.Desc = desc
		}
	}

	if _, err := ctx.FormFile("file"); err == nil {
		if fi, err := service.NewUploader().UploadWithThumbnailSize(ctx, 150, 150); err == nil {
			prodEntity.ImgUrl = fi.OriginUrl
			prodEntity.Thumbx150Url = fi.ThumbUrl
		}
	}

	prodEntity.Update()
	ctx.Redirect(http.StatusSeeOther, "/pai/prod")
}

func prodDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("ID"))
	page, _ := strconv.Atoi(ctx.PostForm("page"))
	if id > 0 {
		_ = service.NewProder().DeleteByID(uint(id))
	}

	const pageSize = 6
	var ps = service.NewProder()
	_, total := ps.GetProdsByPage(page, pageSize)
	targetPage := buildPageInfo(page, pageSize, total).CurrentPage
	if targetPage > 1 {
		ctx.Redirect(http.StatusSeeOther, "/pai/prod?page="+strconv.Itoa(targetPage))
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/pai/prod")
}

// 首页处理函数
func index(ctx *gin.Context) {
	var as = service.NewAbout()
	var ps = service.NewProder()
	var bs = service.NewBanner()

	// 获取数据
	about := as.GetAbout()
	products := ps.GetProds()
	banners := bs.GetBanners()

	// 准备模板数据
	data := gin.H{
		"Title":    "企业门户网站",
		"About":    *about,
		"Products": products,
		"Banners":  banners,
		"Year":     time.Now().Year(),
	}

	ctx.HTML(http.StatusOK, "index.tmpl", data)
}
func about(ctx *gin.Context) {
	var as = service.NewAbout()
	about := as.GetAbout()
	ctx.HTML(http.StatusOK, "about.tmpl", *about)
}

func SetupRouter(r *gin.Engine) *gin.Engine {
	pai := r.Group("/pai")
	{
		pai.GET("/login", login)
		pai.POST("/login", doLogin)
		pai.GET("/home", home)
		pai.GET("/banner", banner)
		pai.POST("/banner/upload", bannerUpload)
		pai.POST("/banner/crop", bannerCrop)
		pai.GET("/brand", brand)
		pai.POST("/brand", brandUpdate)
		pai.GET("/prod", prod)
		pai.POST("/prod/update", prodUpdate)
		pai.POST("/prod/delete", prodDelete)
		pai.POST("/logout", logout)
	}

	// 配置首页路由
	r.GET("/", index)
	r.GET("/about", about)

	return r
}
