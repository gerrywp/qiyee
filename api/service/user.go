package service

import (
	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/common"
)

type UserService struct{}

func NewUser() *UserService {
	us := new(UserService)
	return us
}

func (us *UserService) Login(userName, password string) bool {
	var entity models.User
	entity = entity.GetUser(userName)
	if entity.ID == 0 {
		return false
	} else {
		common.SetCurrent(&entity)
		return entity.Password == password
	}
}
