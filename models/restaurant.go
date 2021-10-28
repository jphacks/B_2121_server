package models

import "github.com/jphacks/B_2121_server/openapi"

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
