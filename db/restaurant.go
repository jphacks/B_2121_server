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

func NewRestaurantRepository(db *sqlx.DB) models.RestaurantRepository {
	return &restaurantRepository{
		db: db,
	}
}

type restaurantRepository struct {
	db *sqlx.DB
}

func (r *restaurantRepository) AddOrUpdateRestaurant(ctx context.Context, restaurant *[]models.SearchApiRestaurant, source models.RestaurantSource) ([]*models.Restaurant, error) {
	ret := make([]*models.Restaurant, 0)

	// TODO: N+1 Problem !
	for _, rest := range *restaurant {
		dbRest, err := models_gen.Restaurants(qm.Where("source = ? AND source_id = ?", source, rest.Id)).One(ctx, r.db)
		if xerrors.Is(err, sql.ErrNoRows) {
			result, err := r.db.ExecContext(ctx, `
INSERT INTO restaurants(name, latitude, longitude, address, url, image_url, source, source_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, rest.Name, rest.Location.Latitude, rest.Location.Longitude, rest.Address, rest.PageUrl, rest.ImageUrl, source, rest.Id)
			if err != nil {
				return nil, xerrors.Errorf("failed to insert into database: %w", err)
			}
			lastInsertId, err := result.LastInsertId()
			if err != nil {
				return nil, xerrors.Errorf("failed to get last insert id: %w", err)
			}
			ret = append(ret, &models.Restaurant{
				Id:       lastInsertId,
				ImageUrl: &rest.ImageUrl,
				Location: rest.Location,
				Name:     rest.Name,
			})
			continue
		}
		if err != nil {
			return nil, xerrors.Errorf("failed to get restaurants from database: %w", err)
		}
		// TODO: Update DB records
		ret = append(ret, fromDbRestaurant(dbRest))
	}

	return ret, nil
}

func (r *restaurantRepository) GetRestaurantById(ctx context.Context, id int64) (*models.Restaurant, error) {
	rest, err := models_gen.FindRestaurant(ctx, r.db, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to get restaurant from database: %w", err)
	}

	// TODO: Update when the information is old
	return fromDbRestaurant(rest), nil
}

func fromDbRestaurant(rest *models_gen.Restaurant) *models.Restaurant {
	return &models.Restaurant{
		Id:       rest.ID,
		ImageUrl: rest.ImageURL.Ptr(),
		Location: models.Location{
			Latitude:  rest.Latitude.Float64, // TODO: Handle Nil ?
			Longitude: rest.Longitude.Float64,
		},
		Name: rest.Name,
	}
}
