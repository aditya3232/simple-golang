package inbound

import "github.com/labstack/echo/v4"

type PingHandlerInterface interface {
	Ping(c echo.Context) error
}
