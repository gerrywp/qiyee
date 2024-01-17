package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type About struct {
	gorm.Model
	Content string
}

func (a *About) Update() {
	DB.Save(a)
}

func (a *About) GetAbout() {
	r := DB.Limit(1).Find(&a)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			fmt.Println("user does not exist!")
		}
		// 数据库链接错误记录日志
		fmt.Println(r.Error)
	}
}
