package router

import (
	"net/http"

	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/api/service"
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
	us := service.New()
	if us.Login(form.UserName) {
		ctx.Redirect(http.StatusOK, "/pai/home")
	} else {
		ctx.Redirect(http.StatusOK, "/pai/login")
	}
}

func home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home2.tmpl", nil)
}

func SetupRouter(r *gin.Engine) *gin.Engine {
	r.GET("/pai/login", login)
	r.POST("/pai/login", doLogin)
	r.GET("/pai/home", home)
	return r
}
