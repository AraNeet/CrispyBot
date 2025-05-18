package server

import (
	"github.com/labstack/echo/v4"
)

const PORT = ":8080"

func StartServer() {
	api := echo.New()
	api.GET("/", func(ctx echo.Context) error {
		return ctx.String(200, "Hello, Echo!")
	})
	api.Logger.Fatal(api.Start(PORT))
}
