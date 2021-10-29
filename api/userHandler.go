package api

import (
	"database/sql"
	"io/ioutil"
	"net/http"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

func (h handler) NewUser(ctx echo.Context) error {
	var req openapi.CreateUserRequest
	err := ctx.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest
	}
	vendor, err := models.FromOpenApiAuthVendor(req.Vendor)
	if err != nil {
		return echo.ErrBadRequest
	}
	user, auth, err := h.userUseCase.NewUser(ctx.Request().Context(), req.Name, vendor)
	if err != nil {
		return err
	}

	userOapi := user.ToOpenApiUser()
	authOapi, err := auth.ToOpenApi()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, openapi.CreateUserResponse{
		AuthInfo: *authOapi,
		User:     *userOapi,
	})
}

func (h handler) GetMyProfile(ctx echo.Context) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}
	userDetail, err := h.userUseCase.MyUser(ctx.Request().Context(), info.UserId)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, userDetail.ToOpenApi())
}

func (h handler) UploadProfileImage(ctx echo.Context) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}
	userId := info.UserId
	data, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	imgUrl, err := h.userUseCase.UpdateUserProfileImage(ctx.Request().Context(), userId, data)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, openapi.UploadImageProfileResponse{ImageUrl: imgUrl})
}

func (h handler) ListUserCommunities(ctx echo.Context, id openapi.Long, params openapi.ListUserCommunitiesParams) error {
	comm, err := h.userUseCase.ListUserCommunities(ctx.Request().Context(), int64(id))
	if err != nil {
		return err
	}

	c := make([]openapi.Community, 0)
	for _, community := range comm {
		c = append(c, *community.ToOpenApiCommunity())
	}
	return ctx.JSON(http.StatusOK, openapi.ListUserCommunityResponse{
		Communities: &c,
		PageInfo:    openapi.PageInfo{},
	})
}

// PostUserMeCommunities - Join community
func (h handler) PostUserMeCommunities(ctx echo.Context) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}
	userId := info.UserId

	var req openapi.PostUserMeCommunitiesJSONRequestBody
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Logger().Errorf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	//TODO: verify invite token
	communityId := 7

	err = h.userUseCase.JoinCommunity(ctx.Request().Context(), userId, int64(communityId))
	if err != nil {
		// TODO: communityがない場合404, Duplicate entryの場合400
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h handler) DeleteUserIdCommunitiesCommunityId(ctx echo.Context, id openapi.Long, communityId openapi.Long) error {
	info := session.GetAuthInfo(ctx)
	if !info.Authenticated {
		return echo.ErrUnauthorized
	}
	userId := info.UserId

	if userId != int64(id) {
		return echo.ErrForbidden
	}

	err := h.userUseCase.LeaveCommunity(ctx.Request().Context(), userId, int64(communityId))
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
