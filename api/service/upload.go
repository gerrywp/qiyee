package service

import (
	"errors"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gerry.wang/qiyee/configs"
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
