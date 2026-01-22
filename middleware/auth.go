package middleware

import (
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	// IgnorePaths 不需要登录检查的路径列表
	IgnorePaths = []string{}
)

// CheckLogin 登录检查中间件
func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查当前路径是否在忽略列表中
		reqPath := path.Clean(ctx.Request.URL.Path)
		for _, p := range IgnorePaths {
			if path.Clean(p) == reqPath {
				ctx.Next()
				return
			}
		}

		// 检查登录状态
		session := sessions.Default(ctx)
		qiyee := session.Get("qiyee")
		if qiyee == nil {
			// 判断是否为AJAX请求
			if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" || ctx.Request.Header.Get("Content-Type") == "application/json" {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code": 0,
					"msg":  "请先登录",
				})
			} else {
				// 当前项目会有问题，因为使用了iframe嵌套，跳转会失败
				ctx.Redirect(http.StatusTemporaryRedirect, "/pai/login")
			}
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
