package models

import (
	"context"

	"github.com/jphacks/B_2121_server/openapi"
)

type RestaurantSource string

const HotpepperSource RestaurantSource = "hotpepper"

type SearchApiRestaurant struct {
	Id       string
	Name     string
	Location Location
	ImageUrl string
	PageUrl  string
	Address  string
}

type Restaurant struct {
	Id       int64
	ImageUrl *string
	Location Location
	Name     string
}

func (r *Restaurant) ToOpenApiRestaurant() *openapi.Restaurant {
	return &openapi.Restaurant{
		Id:       openapi.Long(r.Id),
		ImageUrl: r.ImageUrl,
		Location: r.Location.ToOpenApiLocation(),
		Name:     r.Name,
	}
}

type RestaurantRepository interface {
	AddOrUpdateRestaurant(ctx context.Context, restaurant *[]SearchApiRestaurant, source RestaurantSource) ([]*Restaurant, error)
}
