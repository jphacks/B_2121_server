package models

import "context"

type CommunityRestaurantsRepository interface {
	ListCommunityRestaurants(ctx context.Context, communityId int64) ([]*Restaurant, error)
	AddRestaurants(ctx context.Context, communityId int64, restaurantId int64) error
}
