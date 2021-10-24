package main

import (
	"fmt"

	"github.com/jphacks/B_2121_server/api"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()
	handler := api.NewHandler()
	openapi.RegisterHandlers(e, handler)
	swaggerHandler := echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URL = "/spec"
	})
	e.GET("/spec", api.GetSwaggerSpec)
	e.GET("/swagger/*", swaggerHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", 8080)))
}
