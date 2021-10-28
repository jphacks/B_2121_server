package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
)

func NewRestaurantRepository(db *sqlx.DB) models.RestaurantRepository {
	return &restaurantRepository{
		db: db,
	}
}

type restaurantRepository struct {
	db *sqlx.DB
}
