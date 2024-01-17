package service

import (
	"gerry.wang/qiyee/api/models"
)

type AboutService struct{}

func NewAbout() *AboutService {
	return new(AboutService)
}

func (a *AboutService) Update(content string) {
	entity := new(models.About)
	entity.ID = 1
	entity.Content = content
	entity.Update()
}

func (a *AboutService) GetAbout() *models.About {
	entity := new(models.About)
	entity.GetAbout()
	return entity
}
