package usecase

import (
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/restaurant_search"
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
