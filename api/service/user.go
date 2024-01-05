package service

import (
	"gerry.wang/qiyee/api/models"
	"gerry.wang/qiyee/common"
)

type UserService struct {
	User *models.User
}

func New() *UserService {
	us := new(UserService)
	us.User = new(models.User)
	return us
}

func (us *UserService) Login(userName, password string) bool {
	um := us.User.GetUser(userName)
	if um.ID == 0 {
		return false
	} else {
		common.SetCurrent(&um)
		return um.Password == password
	}
}
