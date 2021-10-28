package usecase

import (
	"context"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/restaurant_search"
	"golang.org/x/xerrors"
)

func NewRestaurantUseCase(restaurantSearch restaurant_search.SearchApi, restaurantRepository models.RestaurantRepository) RestaurantUseCase {
	return RestaurantUseCase{
		restaurantSearch:     restaurantSearch,
		restaurantRepository: restaurantRepository,
	}
}

type RestaurantUseCase struct {
	restaurantSearch     restaurant_search.SearchApi
	restaurantRepository models.RestaurantRepository
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
