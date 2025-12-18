package inbound

import "github.com/labstack/echo/v4"

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
}
