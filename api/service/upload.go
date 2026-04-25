package service

import (
	"errors"
	"image"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gerry.wang/qiyee/configs"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

type UploadService struct{}

func NewUploader() *UploadService {
	us := new(UploadService)
	return us
}

// 上传得文件信息
type FileInfo struct {
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	FileUrl  string `json:"fileUrl"`
	FileType string `json:"fileType"`
}

// BannerUploadResult Banner上传结果（包含原图和缩略图）
type BannerUploadResult struct {
	OriginUrl string `json:"originUrl"` // 原始图片URL
	ThumbUrl  string `json:"thumbUrl"`  // 缩略图URL
}

func (us *UploadService) Upload(ctx *gin.Context) (FileInfo, error) {
	// 获取文件(注意这个地方的file要和html模板中的name一致)
	file, err := ctx.FormFile("file")
	if err != nil {
		return FileInfo{}, errors.New("上传文件不能为空")
	}
	//获取文件的后缀名
	fileExt := path.Ext(file.Filename)

	// 允许上传文件后缀
	allowExt := "jpg,gif,png,bmp,jpeg,JPG"
	// 检查上传文件后缀
	if !checkFileExt(fileExt, allowExt) {
		return FileInfo{}, errors.New("上传文件格式不正确，文件后缀只允许为：" + allowExt + "的文件")
	}
	// 允许文件上传最大值
	allowSize := "5M"
	// 检查上传文件大小
	isvalid, err := checkFileSize(file.Size, allowSize)
	if err != nil {
		return FileInfo{}, err
	}
	if !isvalid {
		return FileInfo{}, errors.New("上传文件大小不得超过：" + allowSize)
	}
	//根据当前时间鹾生成一个新的文件名
	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt, 10)
	//新的文件名
	fileName := fileNameStr + fileExt
	//保存上传文件
	filePath := filepath.Join(configs.GetAbsPath("/web/static/upload"), "/", fileName)
	err2 := ctx.SaveUploadedFile(file, filePath)
	if err2 != nil {
		return FileInfo{}, err2
	}

	// 返回结果
	result := FileInfo{
		FileName: file.Filename,
		FileSize: file.Size,
		FileUrl:  "/static/upload" + "/" + fileName,
	}
	return result, nil
}

// UploadWithThumbnail 上传文件并生成缩略图
func (us *UploadService) UploadWithThumbnail(ctx *gin.Context) (BannerUploadResult, error) {
	// 上传原始文件
	fileInfo, err := us.Upload(ctx)
	if err != nil {
		return BannerUploadResult{}, err
	}

	// 生成缩略图 - 根据banner比例等比例缩放
	originalPath := filepath.Join(configs.GetAbsPath("/web/static/upload"), filepath.Base(fileInfo.FileUrl))
	thumbUrl, err := us.GenerateThumbnail(originalPath, 400, 125) // 400:125 ≈ 1920:600 的缩略图比例
	if err != nil {
		// 缩略图生成失败，但不影响原图
		return BannerUploadResult{
			OriginUrl: fileInfo.FileUrl,
			ThumbUrl:  fileInfo.FileUrl,
		}, nil
	}

	return BannerUploadResult{
		OriginUrl: fileInfo.FileUrl,
		ThumbUrl:  thumbUrl,
	}, nil
}

// GenerateThumbnail 生成缩略图
// originalPath: 原始图片的完整路径
// maxWidth, maxHeight: 缩略图的最大宽高（等比例缩放）
func (us *UploadService) GenerateThumbnail(originalPath string, maxWidth, maxHeight int) (string, error) {
	// 读取原始图片
	src, err := imaging.Open(originalPath)
	if err != nil {
		return "", err
	}

	// 获取原始图片尺寸
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// 计算缩放比例，保持宽高比
	scaleX := float64(maxWidth) / float64(srcWidth)
	scaleY := float64(maxHeight) / float64(srcHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// 计算新的尺寸
	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)

	// 等比例缩放图片
	thumb := imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)

	// 生成缩略图文件名
	fileName := filepath.Base(originalPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	thumbFileName := fileNameWithoutExt + "_thumb" + filepath.Ext(fileName)
	thumbPath := filepath.Join(filepath.Dir(originalPath), thumbFileName)

	// 保存缩略图
	err = imaging.Save(thumb, thumbPath)
	if err != nil {
		return "", err
	}

	// 返回缩略图的URL路径
	thumbUrl := "/static/upload" + "/" + thumbFileName
	return thumbUrl, nil
}

// CropAndSaveImage 根据裁切数据保存裁切后的缩略图（前端上传canvas数据）
// 这个方法用于处理从前端 cropper.js 传过来的裁切数据
func (us *UploadService) CropAndSaveImage(ctx *gin.Context, imageData []byte, width, height, x, y int) (string, error) {
	// 根据裁切数据生成新的图片
	// x, y: 裁切的起始位置
	// width, height: 裁切的宽高

	// 如果前端上传的是 canvas base64 或 blob，则直接保存
	// 如果上传的是坐标数据，则需要从原图中裁切

	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt, 10)
	fileName := fileNameStr + "_thumb.png"
	filePath := filepath.Join(configs.GetAbsPath("/web/static/upload"), fileName)

	// 假设前端上传的是 blob 数据，直接保存即可
	err := os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return "", err
	}

	thumbUrl := "/static/upload" + "/" + fileName
	return thumbUrl, nil
}

// ResizeImage 调整图片大小
func (us *UploadService) ResizeImage(imagePath string, width, height int) (image.Image, error) {
	src, err := imaging.Open(imagePath)
	if err != nil {
		return nil, err
	}
	// 使用 Fit 方法进行等比例缩放
	return imaging.Fit(src, width, height, imaging.Lanczos), nil
}

// 检查文件格式是否合法
func checkFileExt(fileExt string, typeString string) bool {
	// 允许上传文件后缀
	exts := strings.Split(typeString, ",")
	// 是否验证通过
	isValid := false
	for _, v := range exts {
		// 对比文件后缀
		if strings.EqualFold(fileExt, "."+v) {
			isValid = true
			break
		}
	}
	return isValid
}

// 检查上传文件大小
func checkFileSize(fileSize int64, maxSize string) (bool, error) {
	// 匹配上传文件最大值
	match := regexp.MustCompile(`^([0-9]+)(?i:([a-z]*))$`)
	matches := match.FindStringSubmatch(maxSize)
	if len(matches) < 3 {
		err := errors.New("上传文件大小未设置,请在后台配置,格式为(30M,30k,30MB)")
		return false, err
	}
	var cfSize int64
	r, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		r = 0
	}
	switch strings.ToUpper(matches[2]) {
	case "MB", "M":
		cfSize = r * 1024 * 1024
	case "KB", "K":
		cfSize = r * 1024
	case "":
		cfSize = r
	}
	if cfSize == 0 {
		err = errors.New("上传文件大小未设置,请在后台配置,格式为(30M,30k,30MB)最大单位为MB")
		return false, err
	}
	return cfSize >= fileSize, nil
}
