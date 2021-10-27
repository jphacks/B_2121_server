package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"golang.org/x/xerrors"
)

type communityRepository struct {
	db *sqlx.DB
}

func NewCommunityRepository(db *sqlx.DB) models.CommunityRepository {
	return &communityRepository{db}
}

func (c *communityRepository) GetCommunityByID(ctx context.Context, id int64) (*models.CommunityDetail, error) {
	community, err := models_gen.FindCommunity(ctx, c.db, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to find a community by ID: %w", err)
	}

	numRestaurants, err := community.CommunitiesRestaurants().Count(ctx, c.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to count restaurants related to the community: %w", err)
	}

	userCount, err := community.Affiliations().Count(ctx, c.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to count users related to the community: %w", err)
	}

	return &models.CommunityDetail{
		Community:     models.Community{Community: *community},
		NumRestaurant: int(numRestaurants),
		UserCount:     int(userCount),
	}, nil
}
