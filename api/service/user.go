package service

import (
	"encoding/json"

	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserService struct{}

func NewUser() *UserService {
	us := new(UserService)
	return us
}

func (us *UserService) Login(ctx *gin.Context, userName, password string, rememberMe bool) bool {
	var entity models.User
	entity = entity.GetUser(userName)
	if entity.ID == 0 {
		return false
	} else {
		if entity.Password != password {
			return false
		}
		common.SetCurrent(&entity)
		session := sessions.Default(ctx)
		data, _ := json.Marshal(&entity)
		if rememberMe {
			session.Options(sessions.Options{
				MaxAge:   7 * 24 * 60 * 60, //7天
				Path:     "/",
				HttpOnly: true,
			})
		} else {
			session.Options(sessions.Options{
				MaxAge:   60 * 30, //默认30分钟
				Path:     "/",
				HttpOnly: true,
			})
		}
		// 直接使用qiyee这个session以加密字符串形式存储用户对象
		session.Set("qiyee", string(data))
		session.Save()
		return true
	}
}
