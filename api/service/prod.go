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
