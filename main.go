package main

import (
	"fmt"

	"github.com/jphacks/B_2121_server/api"
	"github.com/jphacks/B_2121_server/config"
	database "github.com/jphacks/B_2121_server/db"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/restaurant_search"
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

	db, err := database.New(conf.DBHost, conf.DBPort, conf.DBDatabaseName, conf.DBUser, conf.DBPassword, e.Logger)
	if err != nil {
		e.Logger.Fatalf("failed to connect DB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			e.Logger.Errorf("Failed to close database connection: %v", err)
		}
	}()

	// DB migration
	err = database.Migrate(db, ".")
	if err != nil {
		e.Logger.Errorf("Database migration failed: %v", err)
		return
	}

	hotpepper := restaurant_search.NewSearchApi(conf.HotpepperApiKey)

	// register repositories
	communityRepository := database.NewCommunityRepository(db)
	userRepository := database.NewUserRepository(db)
	affiliationRepository := database.NewAffiliationRepository(db)
	communityRestaurantsRepository := database.NewCommunityRestaurantsRepository(db)
	restaurantRepository := database.NewRestaurantRepository(db)

	store := session.NewStore("key")
	userUseCase := usecase.NewUserUseCase(store, userRepository, conf)
	communityUseCase := usecase.NewCommunityUseCase(store, conf, communityRepository, affiliationRepository, communityRestaurantsRepository)
	restaurantUseCase := usecase.NewRestaurantUseCase(hotpepper, restaurantRepository)
	handler := api.NewHandler(userUseCase, communityUseCase, restaurantUseCase)
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
