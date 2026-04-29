package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/api/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupAdminRouter 配置后端管理路由
func SetupAdminRouter(r *gin.Engine) *gin.Engine {
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
		pai.GET("/site", site)
		pai.POST("/site", siteUpdate)
		pai.GET("/prod", prod)
		pai.POST("/prod/update", prodUpdate)
		pai.POST("/prod/delete", prodDelete)
		pai.GET("/news", news)
		pai.POST("/news/update", newsUpdate)
		pai.POST("/news/delete", newsDelete)
		pai.POST("/logout", logout)
	}
	return r
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

func site(ctx *gin.Context) {
	var ss = service.NewSite()
	siteInfo := ss.GetSite()
	ctx.HTML(http.StatusOK, "site.tmpl", siteInfo)
}

func siteUpdate(ctx *gin.Context) {
	siteInfo := models.Site{
		Model:     gorm.Model{ID: 1},
		Title:     ctx.PostForm("Title"),
		Tel:       ctx.PostForm("Tel"),
		Email:     ctx.PostForm("Email"),
		ICP:       ctx.PostForm("ICP"),
		Copyright: ctx.PostForm("Copyright"),
		FL1:       ctx.PostForm("FL1"),
		FL1URL:    ctx.PostForm("FL1URL"),
		FL2:       ctx.PostForm("FL2"),
		FL2URL:    ctx.PostForm("FL2URL"),
		FL3:       ctx.PostForm("FL3"),
		FL3URL:    ctx.PostForm("FL3URL"),
		FL4:       ctx.PostForm("FL4"),
		FL4URL:    ctx.PostForm("FL4URL"),
		FL5:       ctx.PostForm("FL5"),
		FL5URL:    ctx.PostForm("FL5URL"),
	}
	var ss = service.NewSite()
	ss.Update(siteInfo)
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

func news(ctx *gin.Context) {
	page := 1
	pageParam := ctx.Query("page")
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	const pageSize = 10
	var ns = service.NewNews()
	news, total := ns.GetNewsByPage(page, pageSize)
	ctx.HTML(http.StatusOK, "news.tmpl", gin.H{
		"News": news,
		"Page": buildPageInfo(page, pageSize, total),
	})
}

func newsUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("ID"))
	title := ctx.PostForm("Title")
	content := ctx.PostForm("Content")

	newsEntity := models.News{
		Title:   title,
		Content: content,
	}

	if id > 0 {
		if existing, err := service.NewNews().FindByID(uint(id)); err == nil {
			newsEntity = *existing
			newsEntity.Title = title
			newsEntity.Content = content
		}
	}

	newsEntity.Update()
	ctx.JSON(http.StatusOK, gin.H{"code": true, "msg": "保存成功"})
}

func newsDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("ID"))
	if id > 0 {
		_ = service.NewNews().DeleteByID(uint(id))
	}
	ctx.JSON(http.StatusOK, gin.H{"code": true, "msg": "删除成功"})
}
