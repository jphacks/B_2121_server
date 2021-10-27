package api

import (
	"io/ioutil"
	"net/http"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/session"
	"github.com/labstack/echo/v4"
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
