package service

import (
	"gerry.wang/qiyee/api/models"
)

type SiteService struct{}

func NewSite() *SiteService {
	return new(SiteService)
}

func (s *SiteService) Update(site models.Site) {
	entity := new(models.Site)
	entity.Url = site.Url
	entity.ID = site.ID
	entity.Title = site.Title
	entity.Tel = site.Tel
	entity.Email = site.Email
	entity.ICP = site.ICP
	entity.Copyright = site.Copyright
	entity.FL1 = site.FL1
	entity.FL1URL = site.FL1URL
	entity.FL2 = site.FL2
	entity.FL2URL = site.FL2URL
	entity.FL3 = site.FL3
	entity.FL3URL = site.FL3URL
	entity.FL4 = site.FL4
	entity.FL4URL = site.FL4URL
	entity.FL5 = site.FL5
	entity.FL5URL = site.FL5URL
	entity.Update()
}

func (s *SiteService) GetSite() *models.Site {
	entity := new(models.Site)
	entity.GetSite()
	return entity
}
