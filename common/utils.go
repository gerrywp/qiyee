package common

import (
	"encoding/json"

	"gerry.wang/qiyee/api/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var user models.User

func Current() *models.User {
	if user.ID != 0 {
		return &user
	} else {
		return nil
	}
}

func SetCurrent(u *models.User) {
	user = *u
}

// GenerateSessionID 生成一个唯一的会话ID
func GenerateSessionID() string {
	// 这里可以使用更复杂的生成逻辑，例如UUID
	return "session_" + string(rune(user.ID))
}

// SetCookie 设置一个cookie到客户端浏览器
func SetCookie(ctx *gin.Context, name, value string, maxAge int) {
	// 使用gin.Context设置cookie
	ctx.SetCookie(name, value, maxAge, "/", "", false, true)
	println("Setting cookie:", name, value, "with max age:", maxAge)
}

// StoreSession 存储会话信息
func StoreSession(ctx *gin.Context, sessionID string, user *models.User) {
	// 使用github.com/gin-contrib/sessions库来存储会话信息
	session := sessions.Default(ctx)
	data, _ := json.Marshal(user)
	session.Set(sessionID, string(data))
	if err := session.Save(); err != nil {
		println("Error saving session:", err.Error())
	} else {
		println("Storing session:", sessionID, "for user:", user.UserName)
	}
}
