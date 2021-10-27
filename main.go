package main

import (
	"fmt"

	"github.com/jphacks/B_2121_server/api"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	log.Info("Server starting...")
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	store := session.NewStore("key")
	handler := api.NewHandler(store)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.NewSessionMiddleware(&session.MiddlewareConfig{SessionStore: store}))
	openapi.RegisterHandlers(e, handler)
	swaggerHandler := echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URL = "/spec"
	})
	e.GET("/spec", api.GetSwaggerSpec)
	e.GET("/swagger/*", swaggerHandler)
	e.Static("/images","./profileImages/")
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", 8080)))
}
