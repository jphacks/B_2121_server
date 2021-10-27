package api

import (
	"io/ioutil"
	"net/http"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/session"
	"github.com/jphacks/B_2121_server/usecase"
	"github.com/labstack/echo/v4"
)

func NewHandler(sessionStore session.Store) openapi.ServerInterface {
	return &handler{
		userUseCase: usecase.NewUserUseCase(sessionStore),
	}
}

type handler struct {
	userUseCase usecase.UserUseCase
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
	prof, err := h.userUseCase.UpdateUserProfileImage(userId, data)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, prof)
}

func (h handler) NewCommunity(ctx echo.Context) error {
	panic("implement me")
}

func (h handler) SearchCommunities(ctx echo.Context, params openapi.SearchCommunitiesParams) error {
	panic("implement me")
}

func (h handler) GetCommunityById(ctx echo.Context, id int) error {
	panic("implement me")
}

func (h handler) ListCommunityRestaurants(ctx echo.Context, id int, params openapi.ListCommunityRestaurantsParams) error {
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

func (h handler) ListUsersOfCommunity(ctx echo.Context, id int, params openapi.ListUsersOfCommunityParams) error {
	panic("implement me")
}

func (h handler) SearchRestaurants(ctx echo.Context, params openapi.SearchRestaurantsParams) error {
	panic("implement me")
}

func (h handler) NewUser(ctx echo.Context) error {
	// TODO: Store to Database
	var req openapi.CreateUserRequest
	err := ctx.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest
	}
	vendor, err := models.FromOpenApiAuthVendor(req.Vendor)
	if err != nil {
		return echo.ErrBadRequest
	}
	user, auth, err := h.userUseCase.NewUser(req.Name, vendor)
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
	userDetail, err := h.userUseCase.MyUser(info.UserId)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, userDetail.ToOpenApiUser())
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

func (h handler) ListUserCommunities(ctx echo.Context, id openapi.Long, params openapi.ListUserCommunitiesParams) error {
	panic("implement me")
}
