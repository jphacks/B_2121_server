package models

import "github.com/jphacks/B_2121_server/openapi"

type Location struct {
	Latitude  float64
	Longitude float64
}

func FromOpenApiLocation(location openapi.Location) *Location {
	return &Location{
		Latitude:  location.Lat,
		Longitude: location.Lng,
	}
}
