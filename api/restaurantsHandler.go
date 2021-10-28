package api

import (
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/labstack/echo/v4"
)

func (h handler) GetRestaurantId(ctx echo.Context, id int64) error {
	panic("implement me")
}

func (h handler) AddRestaurantToCommunity(ctx echo.Context, id int) error {
	panic("implement me")
}

func (h handler) RemoveRestaurantFromCommunity(ctx echo.Context, id int64, restaurantId int64) error {
	panic("implement me")
}

func (h handler) GetRestaurantComment(ctx echo.Context, id int, restaurantId int) error {
	panic("implement me")
}

func (h handler) UpdateRestaurantComment(ctx echo.Context, id int, restaurantId int) error {
	panic("implement me")
}

func (h handler) SearchRestaurants(ctx echo.Context, params openapi.SearchRestaurantsParams) error {
	panic("implement me")
}
