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

// Upload 上传Banner（保存原图和缩略图）
func (b *BannerService) Upload(ctx *gin.Context) bool {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	uid := uint(id)
	if fi, err := NewUploader().UploadWithThumbnail(ctx); err != nil {
		return false
	} else {
		entity := new(models.Banner)
		entity.ID = uid
		entity.Url = fi.OriginUrl
		entity.ThumbUrl = fi.ThumbUrl
		entity.Update()
		return true
	}
}

// CropImage 处理图片裁切（接收 base64 图片数据）
func (b *BannerService) CropImage(ctx *gin.Context) (string, error) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	uid := uint(id)

	// 获取裁切后的图片数据（base64 或 blob）
	_ = ctx.PostForm("imageData") // 预留参数

	// 这里简化处理，实际中可能需要更复杂的逻辑
	// 如果是 base64，需要进行解码
	entity := new(models.Banner)
	banners := entity.GetAll()

	for _, banner := range banners {
		if banner.ID == uid {
			// 返回缩略图URL
			return banner.ThumbUrl, nil
		}
	}

	return "", nil
}

func (b *BannerService) GetBanners() []models.Banner {
	entity := new(models.Banner)
	return entity.GetAll()
}
