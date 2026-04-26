package service

import "gerry.wang/qiyee/api/models"

type ProdService struct{}

func NewProder() *ProdService {
	ps := new(ProdService)
	return ps
}

func (ps *ProdService) Update(pm models.Prod) {
	pm.Update()
}

func (ps *ProdService) GetProds() []models.Prod {
	entity := new(models.Prod)
	return entity.GetAll()
}

func (ps *ProdService) GetProdsByPage(page, pageSize int) ([]models.Prod, int64) {
	entity := new(models.Prod)
	return entity.GetPaged(page, pageSize)
}

func (ps *ProdService) FindByID(id uint) (*models.Prod, error) {
	entity := new(models.Prod)
	return entity.FindByID(id)
}

func (ps *ProdService) DeleteByID(id uint) error {
	entity := new(models.Prod)
	return entity.DeleteByID(id)
}
