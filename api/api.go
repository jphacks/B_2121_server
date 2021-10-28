package api

import (
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/usecase"
	"github.com/labstack/echo/v4"
)

func NewHandler(userUseCase usecase.UserUseCase, communityUseCase usecase.CommunityUseCase) openapi.ServerInterface {
	return &handler{
		userUseCase:      userUseCase,
		communityUseCase: communityUseCase,
	}
}

type handler struct {
	userUseCase      usecase.UserUseCase
	communityUseCase usecase.CommunityUseCase
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

func (h handler) GetUserIdBookmark(ctx echo.Context, id openapi.Long) error {
	panic("implement me")
}

func (h handler) PostUserIdBookmark(ctx echo.Context, id openapi.Long) error {
	panic("implement me")
}

func (h handler) DeleteUserIdBookmarkCommunityId(ctx echo.Context, id openapi.Long, communityId openapi.Long) error {
	panic("implement me")
}
