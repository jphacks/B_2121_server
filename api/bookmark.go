package api

import (
	"net/http"

	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/session"
	"github.com/labstack/echo/v4"
)

func (h handler) GetUserIdBookmark(ctx echo.Context, id openapi.Long) error {
	panic("implement me")
}

func (h handler) PostUserIdBookmark(ctx echo.Context, id openapi.Long) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}
	userId := info.UserId

	if int64(id) != userId {
		return echo.ErrUnauthorized
	}

	var req openapi.PostUserIdBookmarkJSONRequestBody
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Logger().Errorf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	err = h.bookmarkUseCase.CreateBookmark(ctx.Request().Context(), userId, int64(req.CommunityId))
	if err != nil {
		// TODO: communityがない場合404, Duplicate entryの場合400
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (h handler) DeleteUserIdBookmarkCommunityId(ctx echo.Context, id openapi.Long, communityId openapi.Long) error {
	panic("implement me")
}
