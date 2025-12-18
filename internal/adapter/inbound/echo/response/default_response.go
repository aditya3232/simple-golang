package response

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type DefaultResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type DefaultResponseWithPaginations struct {
	Message    string      `json:"message"`
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	TotalCount int64 `json:"total_count"`
	Limit      int64 `json:"limit"`
	TotalPage  int64 `json:"total_page"`
}

func ResponseSuccess(message string, data any) DefaultResponse {
	return DefaultResponse{
		Message: message,
		Data:    data,
	}
}

func ResponseSuccessWithPagination(message string, data any, page, totalData, totalPage, limit int64) DefaultResponseWithPaginations {
	return DefaultResponseWithPaginations{
		Message: message,
		Data:    data,
		Pagination: &Pagination{
			Page:       page,
			TotalCount: totalData,
			Limit:      limit,
			TotalPage:  totalPage,
		},
	}
}

func ResponseError(message string) DefaultResponse {
	return DefaultResponse{
		Message: message,
		Data:    nil,
	}
}

func RespondWithError(c echo.Context, code int, context string, err error) error {
	log.Errorf("%s: %v", context, err)
	resp := DefaultResponse{
		Message: err.Error(),
		Data:    nil,
	}
	return c.JSON(code, resp)
}
