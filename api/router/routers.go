package router

import (
	"github.com/gin-gonic/gin"
)

type PageInfo struct {
	CurrentPage int
	PageSize    int
	TotalCount  int64
	TotalPages  int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
	Pages       []int
}

func buildPageInfo(page, pageSize int, totalCount int64) PageInfo {
	if page < 1 {
		page = 1
	}
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	pages := make([]int, 0, totalPages)
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	return PageInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		Pages:       pages,
	}
}

func SetupRouter(r *gin.Engine) *gin.Engine {
	r = SetupFrontendRouter(r)
	r = SetupAdminRouter(r)
	return r
}
