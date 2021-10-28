package models

import (
	"context"

	"github.com/jphacks/B_2121_server/models_gen"
)

type Bookmark struct {
	models_gen.Bookmark
}

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, userId, communityId int64) error
	ListBookmarkByUserId(ctx context.Context, userId int64) ([]Community, error)
	DeleteBookmark(ctx context.Context, userId int64, communityId int64) error
}
