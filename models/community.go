package models

import (
	"context"

	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/jphacks/B_2121_server/openapi"
)

type Community struct {
	models_gen.Community
}

type CommunityDetail struct {
	Community
	NumRestaurant int
	UserCount     int
}

type CommunityRepository interface {
	GetCommunityByID(ctx context.Context, id int64) (*CommunityDetail, error)
}

func (c *CommunityDetail) ToOpenApiCommunityDetail() *openapi.CommunityDetail {
	var loc *openapi.Location
	if c.Latitude.Valid && c.Longitude.Valid {
		loc = &openapi.Location{
			Lat: c.Latitude.Float64,
			Lng: c.Longitude.Float64,
		}
	}

	return &openapi.CommunityDetail{
		Community: openapi.Community{
			Description: &c.Description,
			Id:          openapi.Long(c.ID),
			Location:    loc,
			Name:        c.Name,
		},
		NumRestaurant: &c.NumRestaurant,
		UserCount:     c.UserCount,
	}
}