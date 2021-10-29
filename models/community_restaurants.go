package models

import "context"

type CommunityRestaurantsRepository interface {
	ListCommunityRestaurants(ctx context.Context, communityId int64) ([]*Restaurant, error)
	AddRestaurants(ctx context.Context, communityId int64, restaurantId int64) error
	RemoveRestaurants(ctx context.Context, communityId int64, restaurantId int64) error
	ListCommunitiesWithSameRestaurants(ctx context.Context, restaurantId, communityId int64) ([]*Community, error)
}
