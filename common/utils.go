package common

import "gerry.wang/qiyee/api/models"

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
