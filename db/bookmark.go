package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"golang.org/x/xerrors"
)

type bookmarkRepository struct {
	db *sqlx.DB
}

func (b *bookmarkRepository) ListBookmarkByUserId(ctx context.Context, userId int64) (*[]models.Bookmark, error) {
	//TODO: impl
	return &[]models.Bookmark{
		{Bookmark: models_gen.Bookmark{
			ID:          0,
			CommunityID: 1,
			UserID:      2,
			CreatedAt:   time.Now(),
		}},
	}, nil
}

func (b *bookmarkRepository) CreateBookmark(ctx context.Context, userId, communityId int64) error {
	result, err := b.db.ExecContext(ctx, `INSERT INTO bookmarks(community_id, user_id) VALUES (?, ?)`, communityId, userId)
	if err != nil {
		return xerrors.Errorf("failed to create a bookmark: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return xerrors.Errorf("failed to egt row affected: %w", err)
	}

	if rowsAffected != 1 {
		return xerrors.New("rows affected is not 1")
	}

	return nil
}

func NewBookmarkRepository(db *sqlx.DB) models.BookmarkRepository {
	return &bookmarkRepository{db}
}
