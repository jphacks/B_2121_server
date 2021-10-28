package api

import (
	"database/sql"
	"net/http"

	"github.com/jphacks/B_2121_server/openapi"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

func (h handler) NewCommunity(ctx echo.Context) error {
	panic("implement me")
}

func (h handler) SearchCommunities(ctx echo.Context, params openapi.SearchCommunitiesParams) error {
	panic("implement me")
}

func (h handler) GetCommunityById(ctx echo.Context, id int) error {
	community, err := h.communityUseCase.GetCommunity(ctx.Request().Context(), int64(id))
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}
		return err
	}

	return ctx.JSON(http.StatusOK, community.ToOpenApiCommunityDetail())
}

func (h handler) ListCommunityRestaurants(ctx echo.Context, id int, params openapi.ListCommunityRestaurantsParams) error {
	panic("implement me")
}
