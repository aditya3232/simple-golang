package echo

import (
	"simple-golang/internal/port/inbound"

	"github.com/labstack/echo/v4"
)

func InitRoutes(
	e *echo.Echo,
	mid inbound.MiddlewareAdapterInterface,
	userHandler inbound.UserHandlerInterface,
	pingHandler inbound.PingHandlerInterface,
) {
	e.GET("/ping", pingHandler.Ping)

	e.POST("/signin", userHandler.SignIn)
	e.POST("/signup", userHandler.CreateUserAccount)

	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/users", userHandler.GetUser)
	adminGroup.PUT("/users/:id", userHandler.UpdateUser)
	adminGroup.PUT("/update-password", userHandler.UpdatePassword)
	adminGroup.GET("/users/:id", userHandler.GetUserByID)
	adminGroup.DELETE("/users/:id", userHandler.DeleteUser)
}
