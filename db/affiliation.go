package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"golang.org/x/xerrors"
)

type affiliationRepository struct {
	db *sqlx.DB
}

func NewAffiliationRepository(db *sqlx.DB) models.AffiliationRepository {
	return &affiliationRepository{db}
}

func (a affiliationRepository) JoinCommunity(ctx context.Context, userId, communityId int64) error {
	result, err := a.db.ExecContext(ctx, "INSERT INTO affiliation(community_id, user_id) VALUES (?, ?)", communityId, userId)
	if err != nil {
		return xerrors.Errorf("failed to add affiliation: %w", err)
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
