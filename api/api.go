package api

import (
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/jphacks/B_2121_server/usecase"
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
