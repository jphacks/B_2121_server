package api

import (
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/usecase"
	"github.com/labstack/echo/v4"
)

func NewHandler(
	userUseCase usecase.UserUseCase,
	communityUseCase usecase.CommunityUseCase,
	restaurantUseCase usecase.RestaurantUseCase,
	commentUseCase usecase.CommentUseCase,
	bookmarkUseCase usecase.BookmarkUseCase,
) openapi.ServerInterface {
	return &handler{
		userUseCase:       userUseCase,
		communityUseCase:  communityUseCase,
		commentUseCase:    commentUseCase,
		restaurantUseCase: restaurantUseCase,
		bookmarkUseCase:   bookmarkUseCase,
	}
}

type handler struct {
	userUseCase       usecase.UserUseCase
	communityUseCase  usecase.CommunityUseCase
	restaurantUseCase usecase.RestaurantUseCase
	commentUseCase    usecase.CommentUseCase
	bookmarkUseCase   usecase.BookmarkUseCase
}

func (h handler) DeleteUserIdCommunitiesCommunityId(ctx echo.Context, id openapi.Long, communityId openapi.Long) error {
	panic("implement me")
}
