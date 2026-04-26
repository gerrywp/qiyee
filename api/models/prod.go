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

func (p *Prod) DeleteByID(id uint) error {
	return DB.Delete(&Prod{}, id).Error
}

func (p *Prod) FindByID(id uint) (*Prod, error) {
	var prod Prod
	if err := DB.First(&prod, id).Error; err != nil {
		return nil, err
	}
	return &prod, nil
}

func (p *Prod) GetAll() []Prod {
	var prods []Prod
	DB.Find(&prods)
	return prods
}

func (p *Prod) GetPaged(page, pageSize int) ([]Prod, int64) {
	if page < 1 {
		page = 1
	}
	var total int64
	DB.Model(&Prod{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * pageSize
	var prods []Prod
	DB.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&prods)
	return prods, total
}
