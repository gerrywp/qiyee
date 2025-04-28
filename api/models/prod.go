package models

import (
	"gorm.io/gorm"
)

type Prod struct {
	gorm.Model
	Name         string `binding:"required"`
	Desc         string
	ImgUrl       string
	Thumbx150Url string
}

func (p *Prod) Update() {
	DB.Save(p)
}

func (p *Prod) GetAll() []Prod {
	var prods []Prod
	DB.Find(&prods)
	return prods
}
