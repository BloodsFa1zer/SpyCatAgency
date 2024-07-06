package main

import (
	"github.com/labstack/echo/v4"
	"spyCat/middleware"
	"spyCat/routes"
)

func main() {
	e := echo.New()
	e.HidePort = true

	middleware.UserAuth(e)
	routes.UserRoute(e)

	e.Logger.Fatal(e.Start(":6000"))
}
