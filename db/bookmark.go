package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"golang.org/x/xerrors"
)

type bookmarkRepository struct {
	db *sqlx.DB
}

func (b *bookmarkRepository) ListBookmarkByUserId(ctx context.Context, userId int64) ([]models.Community, error) {
	//TODO: impl
	return []models.Community{}, nil
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
