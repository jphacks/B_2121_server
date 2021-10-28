package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

func NewCommunityRestaurantsRepository(db *sqlx.DB) models.CommunityRestaurantsRepository {
	return &communityRestaurantsRepository{db}
}

type communityRestaurantsRepository struct {
	db *sqlx.DB
}

func (c communityRestaurantsRepository) ListCommunityRestaurants(ctx context.Context, communityId int64) ([]*models.Restaurant, error) {
	restaurants, err := models_gen.Restaurants(
		qm.InnerJoin("communities_restaurants ON restaurants.id = communities_restaurants.restaurant_id"),
		qm.Where("community_id = ?", communityId),
	).All(ctx, c.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to get from database: %w", err)
	}

	ret := make([]*models.Restaurant, 0)
	for _, restaurant := range restaurants {
		ret = append(ret, &models.Restaurant{
			Id:       restaurant.ID,
			ImageUrl: restaurant.ImageURL.Ptr(),
			Location: models.Location{ // TODO: Nullable ??
				Latitude:  restaurant.Latitude.Float64,
				Longitude: restaurant.Longitude.Float64,
			},
			Name: restaurant.Name,
		})
	}
	return ret, nil
}

func (c communityRestaurantsRepository) AddRestaurants(ctx context.Context, communityId int64, restaurantId int64) error {
	count, err := models_gen.CommunitiesRestaurants(qm.Where("community_id = ? AND restaurant_id = ?", communityId, restaurantId)).Count(ctx, c.db)
	if err != nil {
		return xerrors.Errorf("failed to get restaurant from database: %w", err)
	}

	// TODO: Return the dedicated error object
	if count > 0 {
		return echo.ErrBadRequest
	}

	result, err := c.db.ExecContext(ctx, "INSERT INTO communities_restaurants(community_id, restaurant_id) VALUES (?, ?)", communityId, restaurantId)
	if err != nil {
		return xerrors.Errorf("failed to add to database: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return xerrors.Errorf("failed to get # of rows affected: %w", err)
	}
	if rows != 1 {
		return xerrors.New("rows affected is not 1")
	}
	return nil
}

func (c communityRestaurantsRepository) RemoveRestaurants(ctx context.Context, communityId int64, restaurantId int64) error {
	count, err := models_gen.CommunitiesRestaurants(qm.Where("community_id = ? AND restaurant_id = ?", communityId, restaurantId)).DeleteAll(ctx, c.db)
	if err != nil {
		return xerrors.Errorf("failed to get restaurant from database: %w", err)
	}

	// TODO: return a dedicated error
	if count == 0 {
		return echo.ErrNotFound
	}

	return nil
}
