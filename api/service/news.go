package service

import "gerry.wang/qiyee/api/models"

type NewsService struct{}

func NewNews() *NewsService {
	ns := new(NewsService)
	return ns
}

func (ns *NewsService) Update(nm models.News) {
	nm.Update()
}

func (ns *NewsService) GetNews() []models.News {
	entity := new(models.News)
	return entity.GetAll()
}

func (ns *NewsService) GetNewsByPage(page, pageSize int) ([]models.News, int64) {
	entity := new(models.News)
	return entity.GetPaged(page, pageSize)
}

func (ns *NewsService) FindByID(id uint) (*models.News, error) {
	entity := new(models.News)
	return entity.FindByID(id)
}

func (ns *NewsService) DeleteByID(id uint) error {
	entity := new(models.News)
	return entity.DeleteByID(id)
}
