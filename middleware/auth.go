package middleware

import (
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	// IgnorePaths 不需要登录检查的路径列表
	IgnorePaths = []string{}
)

// pathMatches 检查请求路径是否匹配忽略列表中的路径模式
func pathMatches(reqPath string, patterns []string) bool {
	for _, pattern := range patterns {
		// 精确匹配
		if path.Clean(pattern) == reqPath {
			return true
		}

		// 支持Gin风格的路由参数匹配（如 /news/:id）
		routePattern := convertGinPatternToRegex(pattern)
		if regex, err := regexp.Compile(routePattern); err == nil {
			if regex.MatchString(reqPath) {
				return true
			}
		}
	}
	return false
}

// convertGinPatternToRegex 将Gin路由模式转换为正则表达式
// 例如：/news/:id -> /news/[0-9]+
// 例如：/user/:id/post/:postId -> /user/[0-9]+/post/[0-9]+
func convertGinPatternToRegex(pattern string) string {
	// 替换 :参数名 为 [0-9]+（匹配数字）或 [^/]+ （匹配非/的任意字符）
	pattern = strings.ReplaceAll(pattern, ":id", "[0-9]+")
	pattern = strings.ReplaceAll(pattern, ":ID", "[0-9]+")
	pattern = strings.ReplaceAll(pattern, ":postId", "[0-9]+")
	pattern = strings.ReplaceAll(pattern, ":newsId", "[0-9]+")
	pattern = strings.ReplaceAll(pattern, ":productId", "[0-9]+")
	// 对于其他可能的参数名，使用更通用的匹配
	regexStr := regexp.MustCompile(`:[a-zA-Z]+`).ReplaceAllString(pattern, `[^/]+`)
	// 转义特殊字符（除了已经是正则的部分）
	regexStr = `^` + regexStr + `$`
	return regexStr
}

// CheckLogin 登录检查中间件
func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查当前路径是否在忽略列表中
		reqPath := path.Clean(ctx.Request.URL.Path)
		if pathMatches(reqPath, IgnorePaths) {
			ctx.Next()
			return
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
