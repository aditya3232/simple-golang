package inbound

import "github.com/labstack/echo/v4"

type UserHandlerInterface interface {
	SignIn(c echo.Context) error
	CreateUserAccount(c echo.Context) error
	UpdatePassword(c echo.Context) error

	GetUser(c echo.Context) error
	GetUserByID(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}
