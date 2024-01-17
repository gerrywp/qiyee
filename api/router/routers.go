package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/api/service"
	"gerry.wang/qiyee/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
	if us.Login(form.UserName, form.Password) {
		session := sessions.Default(ctx)
		if !form.IsRemember {
			session.Options(sessions.Options{
				MaxAge: 0,
			})
		}
		if cu := common.Current(); cu != nil {
			data, _ := json.Marshal(&cu)
			session.Set("user", string(data))
			if err := session.Save(); err != nil {
				fmt.Println(err)
			}
		}
		//ctx.Redirect(http.StatusFound, "/pai/home")
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "登录成功"})
	} else {
		//ctx.Redirect(http.StatusNotModified, "/pai/login")
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "用户名密码错误!"})
	}
}

func home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if r := session.Get("user"); r == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/pai/login")
	} else {
		var cu models.User
		sr := r.(string)
		json.Unmarshal([]byte(sr), &cu)
		ctx.HTML(http.StatusOK, "ihome.tmpl", nil)
	}
}

func banner(ctx *gin.Context) {
	var bs = service.NewBanner()
	banners := bs.GetBanners()
	m := make(map[string]string)
	if len(banners) >= 3 {
		m["url1"] = banners[0].Url
		m["url2"] = banners[1].Url
		m["url3"] = banners[2].Url
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
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "上传成功"})
}

func SetupRouter(r *gin.Engine) *gin.Engine {
	pai := r.Group("/pai")
	{
		pai.GET("/login", login)
		pai.POST("/login", doLogin)
		pai.GET("/home", home)
		pai.GET("/banner", banner)
		pai.POST("/banner/upload", bannerUpload)
		pai.GET("/brand", brand)
		pai.POST("/brand", brandUpdate)
	}
	return r
}
