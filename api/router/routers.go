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
	us := service.New()
	if us.Login(form.UserName, form.Password) {
		session := sessions.Default(ctx)
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
	ctx.HTML(http.StatusOK, "banner.tmpl", nil)
}

func SetupRouter(r *gin.Engine) *gin.Engine {
	r.GET("/pai/login", login)
	r.POST("/pai/login", doLogin)
	r.GET("/pai/home", home)
	r.GET("/pai/banner", banner)
	return r
}
