package models

import (
	"log"

	"gorm.io/gorm"
)

type Banner struct {
	gorm.Model
	Url string
}

func (b *Banner) Insert() uint {
	result := DB.Create(b)
	if result.Error != nil {
		log.Default().Println(result.Error.Error())
		return 0
	}
	return b.ID
}

func (b *Banner) Update() {
	//DB.Model(b).Updates(Banner{Url: b.Url, Model: gorm.Model{UpdatedAt: b.UpdatedAt}})
	DB.Save(b)
}

func (b *Banner) GetAll() []Banner {
	var banners []Banner
	DB.Find(&banners)
	return banners
}
