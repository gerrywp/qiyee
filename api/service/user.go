package service

import (
	"gerry.wang/qiyee/api/models"
)

type UserService struct {
	User *models.User
}

func New() *UserService {
	us := new(UserService)
	us.User = new(models.User)
	return us
}

func (us *UserService) Login(userName string) bool {
	us.User.GetUser(userName)
	return false
}
