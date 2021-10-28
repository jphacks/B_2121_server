package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
)

type bookmarkRepository struct {
	db *sqlx.DB
}

func (b *bookmarkRepository) CreateBookmark(ctx context.Context, userId, communityId int64) (*models.Bookmark, error) {
	//	TODO: impl
	return &models.Bookmark{Bookmark: models_gen.Bookmark{
		ID:          0,
		CommunityID: 1,
		UserID:      2,
		CreatedAt:   time.Now(),
	}}, nil
}

func NewBookmarkRepository(db *sqlx.DB) models.BookmarkRepository {
	return &bookmarkRepository{db}
}
