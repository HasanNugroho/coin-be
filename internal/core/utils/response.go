package utils

type Response struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Meta       *Pagination `json:"meta,omitempty"`
}

type Pagination struct {
	Page        int64 `json:"page"`
	Limit       int64 `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int64 `json:"totalPages"`
	HasPrevPage bool  `json:"hasPrevPage"`
	HasNextPage bool  `json:"hasNextPage"`
}

func NewResponse(statusCode int, message string, data interface{}) *Response {
	return &Response{
		Success:    statusCode >= 200 && statusCode < 300,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func NewResponseWithPagination(statusCode int, message string, data interface{}, pagination *Pagination) *Response {
	return &Response{
		Success:    statusCode >= 200 && statusCode < 300,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       pagination,
	}
}

func NewSuccessResponse(message string, data interface{}) *Response {
	return &Response{
		Success:    true,
		StatusCode: 200,
		Message:    message,
		Data:       data,
	}
}

func NewSuccessResponseWithPagination(message string, data interface{}, pagination *Pagination) *Response {
	return &Response{
		Success:    true,
		StatusCode: 200,
		Message:    message,
		Data:       data,
		Meta:       pagination,
	}
}

func NewCreatedResponse(message string, data interface{}) *Response {
	return &Response{
		Success:    true,
		StatusCode: 201,
		Message:    message,
		Data:       data,
	}
}

func NewErrorResponse(statusCode int, message string) *Response {
	return &Response{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	}
}

func CalculatePagination(page, limit, total int64) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	return &Pagination{
		Page:        page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
		HasPrevPage: page > 1,
		HasNextPage: page < totalPages,
	}
}
