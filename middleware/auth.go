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
				// 更可靠的方法是返回一个HTML页面，其中包含JavaScript来强制在顶层窗口跳转
				redirectHTML := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>重定向</title>
    <script type="text/javascript">
        // 强制在顶层窗口跳转
        if (window.top !== window.self) {
            window.top.location.href = "/pai/login";
        } else {
            window.location.href = "/pai/login";
        }
    </script>
</head>
<body>
    <p>正在跳转登录页面...</p>
</body>
</html>
`
				ctx.Header("Content-Type", "text/html; charset=utf-8")
				ctx.String(http.StatusOK, redirectHTML)
			}
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
