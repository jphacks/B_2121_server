package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

func NewCommunityRestaurantsRepository(db *sqlx.DB) models.CommunityRestaurantsRepository {
	return &communityRestaurantsRepository{db}
}

type communityRestaurantsRepository struct {
	db *sqlx.DB
}

func (c communityRestaurantsRepository) ListCommunityRestaurants(ctx context.Context, communityId int64) ([]*models.Restaurant, error) {
	restaurants, err := models_gen.Restaurants(
		qm.InnerJoin("communities_restaurants ON restaurants.id = communities_restaurants.restaurant_id"),
		qm.Where("community_id = ?", communityId),
	).All(ctx, c.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to get from database: %w", err)
	}

	ret := make([]*models.Restaurant, 0)
	for _, restaurant := range restaurants {
		ret = append(ret, &models.Restaurant{
			Id:       restaurant.ID,
			ImageUrl: restaurant.ImageURL.Ptr(),
			Location: models.Location{ // TODO: Nullable ??
				Latitude:  restaurant.Latitude.Float64,
				Longitude: restaurant.Longitude.Float64,
			},
			Name: restaurant.Name,
		})
	}
	return ret, nil
}

func (c communityRestaurantsRepository) AddRestaurants(ctx context.Context, communityId int64, restaurantId int64) error {
	count, err := models_gen.CommunitiesRestaurants(qm.Where("community_id = ? AND restaurant_id = ?", communityId, restaurantId)).Count(ctx, c.db)
	if err != nil {
		return xerrors.Errorf("failed to get restaurant from database: %w", err)
	}

	// TODO: Return the dedicated error object
	if count > 0 {
		return echo.ErrBadRequest
	}

	result, err := c.db.ExecContext(ctx, "INSERT INTO communities_restaurants(community_id, restaurant_id) VALUES (?, ?)", communityId, restaurantId)
	if err != nil {
		return xerrors.Errorf("failed to add to database: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return xerrors.Errorf("failed to get # of rows affected: %w", err)
	}
	if rows != 1 {
		return xerrors.New("rows affected is not 1")
	}
	return nil
}

func (c communityRestaurantsRepository) RemoveRestaurants(ctx context.Context, communityId int64, restaurantId int64) error {
	count, err := models_gen.CommunitiesRestaurants(qm.Where("community_id = ? AND restaurant_id = ?", communityId, restaurantId)).DeleteAll(ctx, c.db)
	if err != nil {
		return xerrors.Errorf("failed to get restaurant from database: %w", err)
	}

	// TODO: return a dedicated error
	if count == 0 {
		return echo.ErrNotFound
	}

	return nil
}

func (c communityRestaurantsRepository) ListCommunitiesWithSameRestaurants(ctx context.Context, restaurantId, communityId int64) ([]*models.Community, error) {
	communities, err := models_gen.Communities(
		qm.InnerJoin("communities_restaurants ON community_id = communities.id"),
		qm.Where("restaurant_id = ? AND community_id <> ?", restaurantId, communityId),
		qm.OrderBy("communities_restaurants.created_at DESC"),
	).All(ctx, c.db)
	if xerrors.Is(err, sql.ErrNoRows) {
		return []*models.Community{}, nil
	}
	if err != nil {
		return nil, xerrors.Errorf("failed to get community from database: %w", err)
	}

	images := make([]communityImage, 0)
	err = c.db.SelectContext(ctx, &images, `SELECT joined.community_id AS community_id, image_url
FROM (
       SELECT community_id,
              r.image_url,
              ROW_NUMBER() OVER (PARTITION BY community_id ORDER BY r.created_at DESC ) AS num
       FROM communities_restaurants
              INNER JOIN restaurants r
                         ON communities_restaurants.restaurant_id = r.id
       ORDER BY community_id, num) AS joined
       INNER JOIN communities_restaurants ON joined.community_id = communities_restaurants.community_id
WHERE joined.num <= ?
  AND communities_restaurants.restaurant_id = ?
  AND communities_restaurants.community_id <> ?;
`, numOfRestaurantImagePerCommunity, restaurantId, communityId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get image urls from database: %w", err)
	}
	imageUrlMap := toImageUrlMap(images)

	restaurantNums := make([]communityNums, 0, len(communities))
	err = c.db.SelectContext(ctx, &restaurantNums, `SELECT c.community_id AS community_id, COUNT(communities_restaurants.restaurant_id) AS num
FROM (SELECT community_id
      FROM communities_restaurants
      WHERE restaurant_id = ?
        AND community_id <> ?
     ) AS c
       INNER JOIN communities_restaurants ON c.community_id = communities_restaurants.community_id
GROUP BY c.community_id`, restaurantId, communityId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get restaurant count from database: %s", err)
	}
	restaurantNumsMap := toCommunityNumsMap(restaurantNums)

	userNums := make([]communityNums, 0, len(communities))
	err = c.db.SelectContext(ctx, &userNums, `SELECT c.id AS community_id, COUNT(a.user_id) AS num
FROM communities_restaurants
       INNER JOIN communities c ON communities_restaurants.community_id = c.id
       INNER JOIN affiliation a ON c.id = a.community_id
WHERE restaurant_id = ?
  AND communities_restaurants.community_id <> ?
GROUP BY a.community_id`, restaurantId, communityId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get user count from database: %s", err)
	}
	userNumsMap := toCommunityNumsMap(userNums)

	ret := make([]*models.Community, 0, len(communities))
	for _, c := range communities {
		images, ok := imageUrlMap[c.ID]
		if !ok {
			images = []string{}
		}
		restaurantNum, ok := restaurantNumsMap[c.ID]
		if !ok {
			restaurantNum = 0
		}

		userNum, ok := userNumsMap[c.ID]
		if !ok {
			userNum = 0
		}
		ret = append(ret, &models.Community{
			Community:      *c,
			ImageUrls:      images,
			NumRestaurants: restaurantNum,
			NumUsers:       userNum,
		})
	}
	return ret, nil
}
