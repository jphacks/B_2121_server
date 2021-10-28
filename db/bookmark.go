package db

import (
	"context"
	"database/sql"

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
		ret = append(ret, models.Community{Community: *comm})
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

func (b *bookmarkRepository) DeleteBookmark(ctx context.Context, userId int64, communityId int64) error {
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	bm, err := models_gen.Bookmarks(qm.Where("user_id = ?", userId), qm.And("community_id = ?", communityId)).One(ctx, tx)
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return xerrors.Errorf("affiliation not found: %w", err)
		}
		return err
	}

	rowsAffected, err := bm.Delete(ctx, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return xerrors.Errorf("failed to get row affected: %w", err)
	}
	if rowsAffected != 1 {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return xerrors.New("rows affected is not 1")
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func NewBookmarkRepository(db *sqlx.DB) models.BookmarkRepository {
	return &bookmarkRepository{db}
}
