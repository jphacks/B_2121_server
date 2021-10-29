package usecase

import (
	"context"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/restaurant_search"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

func NewRestaurantUseCase(
	restaurantSearch restaurant_search.SearchApi, restaurantRepository models.RestaurantRepository,
	userRepository models.UserRepository, communityRestaurantsRepository models.CommunityRestaurantsRepository) RestaurantUseCase {
	return RestaurantUseCase{
		restaurantSearch:               restaurantSearch,
		restaurantRepository:           restaurantRepository,
		userRepository:                 userRepository,
		communityRestaurantsRepository: communityRestaurantsRepository,
	}
}

type RestaurantUseCase struct {
	restaurantSearch               restaurant_search.SearchApi
	restaurantRepository           models.RestaurantRepository
	userRepository                 models.UserRepository
	communityRestaurantsRepository models.CommunityRestaurantsRepository
}

func (r RestaurantUseCase) SearchRestaurant(ctx context.Context, keyword string, loc *models.Location) ([]*models.Restaurant, error) {
	rest, err := r.restaurantSearch.Search(keyword, loc, 50)
	if err != nil {
		return nil, xerrors.Errorf("failed to search restaurants: %w", err)
	}
	modelRest, err := r.restaurantRepository.AddOrUpdateRestaurant(ctx, rest, r.restaurantSearch.Source())
	if err != nil {
		return nil, xerrors.Errorf("failed to insert or update restaurants in database: %w", err)
	}
	return modelRest, err
}

func (r RestaurantUseCase) GetRestaurantById(ctx context.Context, id int64) (*models.Restaurant, error) {
	rest, err := r.restaurantRepository.GetRestaurantById(ctx, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to get restaurant: %w", err)
	}
	return rest, nil
}

func (r RestaurantUseCase) AddRestaurantToCommunity(ctx context.Context, userId int64, communityId int64, restaurantId int64) error {
	exist, err := r.userRepository.ExistInCommunity(ctx, userId, communityId)
	if err != nil {
		return xerrors.Errorf("failed to whether the user belongs to the community: %w", err)
	}

	if !exist {
		return echo.ErrForbidden
	}

	err = r.communityRestaurantsRepository.AddRestaurants(ctx, communityId, restaurantId)
	if err != nil {
		return xerrors.Errorf("failed to add restaurant: %w", err)
	}

	return nil
}

func (r RestaurantUseCase) RemoveRestaurantFromCommunity(ctx context.Context, userId int64, communityId int64, restaurantId int64) error {
	exist, err := r.userRepository.ExistInCommunity(ctx, userId, communityId)
	if err != nil {
		return xerrors.Errorf("failed to whether the user belongs to the community: %w", err)
	}

	if !exist {
		return echo.ErrForbidden
	}

	err = r.communityRestaurantsRepository.RemoveRestaurants(ctx, communityId, restaurantId)
	if err != nil {
		return xerrors.Errorf("failed to add restaurant: %w", err)
	}

	return nil
}

func (r RestaurantUseCase) GetOtherCommunitiesWithSameRestaurants(ctx context.Context, restaurantId, communityId int64) ([]*models.Community, error) {
	communities, err := r.communityRestaurantsRepository.ListCommunitiesWithSameRestaurants(ctx, restaurantId, communityId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get communities with the same restaurants: %w", err)
	}
	return communities, nil
}
