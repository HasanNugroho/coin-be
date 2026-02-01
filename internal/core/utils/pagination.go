package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int64
	PageSize int64
}

// SortParams holds sorting parameters
type SortParams struct {
	SortBy    string // field name to sort by
	SortOrder string // "asc" or "desc"
}

// PaginationMeta holds pagination metadata for response
type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
}

// ParsePaginationParams extracts and validates pagination parameters from query
func ParsePaginationParams(ctx *gin.Context, defaultPageSize int64) PaginationParams {
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", strconv.FormatInt(defaultPageSize, 10))

	var pageNum int64 = 1
	var pageSizeNum int64 = defaultPageSize

	if p, err := strconv.ParseInt(page, 10, 64); err == nil && p > 0 {
		pageNum = p
	}

	if ps, err := strconv.ParseInt(pageSize, 10, 64); err == nil && ps > 0 {
		pageSizeNum = ps
	}

	// Cap page size at maximum
	if pageSizeNum > 100 {
		pageSizeNum = 100
	}

	return PaginationParams{
		Page:     pageNum,
		PageSize: pageSizeNum,
	}
}

// ParseSortParams extracts and validates sorting parameters from query
func ParseSortParams(ctx *gin.Context, allowedFields []string, defaultField string) SortParams {
	sortBy := ctx.DefaultQuery("sort_by", defaultField)
	sortOrder := ctx.DefaultQuery("sort_order", "desc")

	// Validate sort field
	isAllowed := false
	for _, field := range allowedFields {
		if field == sortBy {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		sortBy = defaultField
	}

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	return SortParams{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(total int64, page int64, pageSize int64) PaginationMeta {
	totalPages := (total + pageSize - 1) / pageSize
	return PaginationMeta{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// BuildPaginatedResponse builds a paginated response with data and metadata
func BuildPaginatedResponse(data interface{}, meta PaginationMeta) map[string]interface{} {
	return map[string]interface{}{
		"data": data,
		"meta": meta,
	}
}
