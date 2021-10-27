package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"golang.org/x/xerrors"
)

type CommunityRepository struct {
	db *sqlx.DB
}

func NewCommunityRepository(db *sqlx.DB) models.CommunityRepository {
	return &CommunityRepository{db}
}

func (c *CommunityRepository) GetCommunityByID(ctx context.Context, id int64) (*models_gen.Community, error) {
	community, err := models_gen.FindCommunity(ctx, c.db, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to find a community by ID: %w", err)
	}

	return community, nil
}
