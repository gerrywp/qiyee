package models

import (
	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title   string `binding:"required"`
	Content string
}

func (n *News) Update() {
	DB.Save(n)
}

func (n *News) DeleteByID(id uint) error {
	return DB.Delete(&News{}, id).Error
}

func (n *News) FindByID(id uint) (*News, error) {
	var news News
	if err := DB.First(&news, id).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (n *News) GetAll() []News {
	var news []News
	DB.Find(&news)
	return news
}

func (n *News) GetPaged(page, pageSize int) ([]News, int64) {
	if page < 1 {
		page = 1
	}
	var total int64
	DB.Model(&News{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * pageSize
	var news []News
	DB.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&news)
	return news, total
}
