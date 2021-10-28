package api

import (
	"database/sql"
	"net/http"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

func (h handler) GetRestaurantId(ctx echo.Context, id int64) error {
	rest, err := h.restaurantUseCase.GetRestaurantById(ctx.Request().Context(), id)
	if xerrors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, rest.ToOpenApiRestaurant())
}

func (h handler) AddRestaurantToCommunity(ctx echo.Context, id int) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}

	var req openapi.AddRestaurantRequest
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Logger().Errorf("failed to bind request body: %v", err)
		return echo.ErrBadRequest
	}

	err = h.restaurantUseCase.AddRestaurantToCommunity(ctx.Request().Context(), info.UserId, int64(id), int64(req.RestaurantId))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h handler) RemoveRestaurantFromCommunity(ctx echo.Context, id int64, restaurantId int64) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}

	err := h.restaurantUseCase.RemoveRestaurantFromCommunity(ctx.Request().Context(), info.UserId, id, restaurantId)

	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (h handler) GetRestaurantComment(ctx echo.Context, id int, restaurantId int) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}

	comm, err := h.commentUseCase.GetComment(ctx.Request().Context(), info.UserId, int64(id), int64(restaurantId))
	if err != nil {
		return err
	}
	comId := openapi.Long(id)
	restId := openapi.Long(restaurantId)

	return ctx.JSON(http.StatusOK, openapi.Comment{
		Body:         &comm,
		CommunityId:  &comId,
		RestaurantId: &restId,
	})
}

func (h handler) UpdateRestaurantComment(ctx echo.Context, id int, restaurantId int) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}

	var req openapi.UpdateRestaurantCommentJSONRequestBody
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Logger().Errorf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	err = h.commentUseCase.SetComment(ctx.Request().Context(), info.UserId, int64(id), int64(restaurantId), *req.Body)
	if err != nil {
		return err
	}
	comId := openapi.Long(id)
	restId := openapi.Long(restaurantId)

	return ctx.JSON(http.StatusOK, openapi.Comment{
		Body:         req.Body,
		CommunityId:  &comId,
		RestaurantId: &restId,
	})
}

func (h handler) SearchRestaurants(ctx echo.Context, params openapi.SearchRestaurantsParams) error {
	var center *models.Location = nil

	ret := make([]openapi.Restaurant, 0)
	if params.Center != nil && params.Center.Lng != 0 && params.Center.Lat != 0 {
		center = models.FromOpenApiLocation(*params.Center)
	}
	rest, err := h.restaurantUseCase.SearchRestaurant(ctx.Request().Context(), params.Keyword, center)
	if xerrors.Is(err, sql.ErrNoRows) {
		return ctx.JSON(http.StatusOK, openapi.SearchRestaurantResponse{
			Restaurants: &ret,
		})
	}
	if err != nil {
		return err
	}
	for _, restaurant := range rest {
		r := restaurant.ToOpenApiRestaurant()
		ret = append(ret, *r)
	}

	return ctx.JSON(http.StatusOK, openapi.SearchRestaurantResponse{
		Restaurants: &ret,
	})
}
