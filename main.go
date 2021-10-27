package main

import (
	"fmt"

	"github.com/jphacks/B_2121_server/api"
	"github.com/jphacks/B_2121_server/config"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/session"
	"github.com/jphacks/B_2121_server/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	log.Info("Server starting...")
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	conf, err := config.LoadServerConfig()
	if err != nil {
		e.Logger.Fatalf("failed to load configuration %v", err)
	}

	store := session.NewStore("key")
	userUseCase := usecase.NewUserUseCase(store, conf)
	handler := api.NewHandler(userUseCase)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.NewSessionMiddleware(&session.MiddlewareConfig{SessionStore: store}))
	openapi.RegisterHandlers(e, handler)
	swaggerHandler := echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URL = "/spec"
	})
	e.GET("/spec", api.GetSwaggerSpec)
	e.GET("/swagger/*", swaggerHandler)
	e.Static("/images", "./profileImages/")
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", 8080)))
}
