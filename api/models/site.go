package models

import (
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	Title     string
	Tel       string
	Email     string
	ICP       string // 备案号
	Copyright string
	FL1       string // 友情链接1
	FL1URL    string
	FL2       string // 友情链接2
	FL2URL    string
	FL3       string // 友情链接3
	FL3URL    string
	FL4       string
	FL4URL    string
	FL5       string
	FL5URL    string
}

func (s *Site) Update() {
	DB.Save(s)
}

func (s *Site) GetSite() {
	r := DB.Limit(1).Find(&s)
	if r.Error != nil {
		// 忽略空记录或查询异常，页面仍然可以显示空表单
		return
	}
}
