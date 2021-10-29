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

func (h handler) NewCommunity(ctx echo.Context) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}

	var req openapi.CreateCommunityRequest
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Logger().Errorf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	loc := models.FromOpenApiLocation(req.Location)
	community, err := h.communityUseCase.NewCommunity(ctx.Request().Context(), info.UserId, req.Name, req.Description, *loc)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, community.ToOpenApiCommunity())
}

func (h handler) SearchCommunities(ctx echo.Context, params openapi.SearchCommunitiesParams) error {
	comm, err := h.communityUseCase.SearchCommunity(ctx.Request().Context(), params.Keyword)
	if err != nil {
		return err
	}

	// TODO: Support search based on location
	ret := make([]openapi.Community, 0)
	for _, community := range comm {
		ret = append(ret, *community.ToOpenApiCommunity())
	}
	return ctx.JSON(http.StatusOK, openapi.SearchCommunityResponse{
		Communities: &ret,
	})
}

func (h handler) GetCommunityById(ctx echo.Context, id int) error {
	community, err := h.communityUseCase.GetCommunity(ctx.Request().Context(), int64(id))
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}
		return err
	}

	return ctx.JSON(http.StatusOK, community.ToOpenApiCommunity())
}

func (h handler) ListCommunityRestaurants(ctx echo.Context, id int, _ openapi.ListCommunityRestaurantsParams) error {
	rest, err := h.communityUseCase.ListRestaurants(ctx.Request().Context(), int64(id))
	if err != nil {
		return err
	}

	ret := make([]openapi.Restaurant, 0)
	for _, restaurant := range rest {
		ret = append(ret, *restaurant.ToOpenApiRestaurant())
	}

	return ctx.JSON(http.StatusOK, openapi.ListCommunityRestaurantsResponse{
		Restaurants: &ret,
	})
}

func (h handler) ListUsersOfCommunity(ctx echo.Context, id int, _ openapi.ListUsersOfCommunityParams) error {
	users, err := h.communityUseCase.ListUsers(ctx.Request().Context(), int64(id))
	if xerrors.Is(err, sql.ErrNoRows) {
		return ctx.JSON(http.StatusOK, openapi.ListCommunityUsersResponse{
			User:  openapi.User{},
			Users: nil,
		})
	}
	if err != nil {
		return err
	}

	openapiUsers := make([]openapi.User, 0)
	for _, u := range users {
		openapiUsers = append(openapiUsers, *u.ToOpenApiUser())
	}
	return ctx.JSON(http.StatusOK, openapi.ListCommunityUsersResponse{
		Users: &openapiUsers,
	})
}

func (h handler) GetCommunityIdToken(ctx echo.Context, id int) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}
	userId := info.UserId

	inviteToken, err := h.communityUseCase.IssueInviteToken(ctx.Request().Context(), userId, int64(id))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.GetCommunityIdTokenResponse{
		ExpiresIn:   int(inviteToken.ExpiresIn.Seconds()),
		InviteToken: inviteToken.Token,
	})
}

func (h handler) UpdateCommunity(ctx echo.Context, id int64) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}
	userId := info.UserId

	var req openapi.UpdateCommunityJSONRequestBody
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Logger().Errorf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	community, err := h.communityUseCase.UpdateCommunity(ctx.Request().Context(), userId, id, req.Name, req.Description, *models.FromOpenApiLocation(req.Location))
	if xerrors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, community.ToOpenApiCommunity())
}
