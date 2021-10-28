package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

type bookmarkRepository struct {
	db *sqlx.DB
}

func (b *bookmarkRepository) ListBookmarkByUserId(ctx context.Context, userId int64) ([]models.Community, error) {
	comms, err := models_gen.Communities(
		qm.InnerJoin("bookmarks ON bookmarks.community_id = communities.id"),
		qm.Where("user_id = ?", userId),
	).All(ctx, b.db)
	if err != nil {
		return nil, err
	}

	ret := make([]models.Community, 0, len(comms))
	for _, comm := range comms {
		ret = append(ret, models.Community{*comm})
	}

	return ret, nil
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
