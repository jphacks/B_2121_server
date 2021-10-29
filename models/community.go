package models

import (
	"context"

	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/jphacks/B_2121_server/openapi"
)

type Community struct {
	models_gen.Community
	ImageUrls      []string
	NumRestaurants int
	NumUsers       int
}

type CommunityRepository interface {
	GetCommunityByID(ctx context.Context, id int64) (*Community, error)
	NewCommunity(ctx context.Context, name string, description string, loc Location) (*Community, error)
	SearchCommunity(ctx context.Context, keyword string) ([]*Community, error)
	UpdateCommunity(ctx context.Context, communityId int64, name string, description string, loc Location) (*Community, error)
}

func (c *Community) ToOpenApiCommunity() *openapi.Community {
	var loc openapi.Location
	if c.Latitude.Valid && c.Longitude.Valid {
		loc = openapi.Location{
			Lat: c.Latitude.Float64,
			Lng: c.Longitude.Float64,
		}
	}
	return &openapi.Community{
		Description:   c.Description,
		Id:            openapi.Long(c.ID),
		Location:      loc,
		Name:          c.Name,
		ImageUrls:     c.ImageUrls,
		NumRestaurant: c.NumRestaurants,
		NumUser:       c.NumUsers,
	}
}
