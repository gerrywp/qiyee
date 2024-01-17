package service

import (
	"strconv"

	"gerry.wang/qiyee/api/models"
	"github.com/gin-gonic/gin"
)

type BannerService struct{}

func NewBanner() *BannerService {
	return new(BannerService)
}

func (b *BannerService) Upload(ctx *gin.Context) bool {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	uid := uint(id)
	if fi, err := NewUploader().Upload(ctx); err != nil {
		return false
	} else {
		entity := new(models.Banner)
		entity.ID = uid
		entity.Url = fi.FileUrl
		entity.Update()
		return true
	}
}

func (b *BannerService) GetBanners() []models.Banner {
	entity := new(models.Banner)
	return entity.GetAll()
}
