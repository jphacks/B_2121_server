package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

const numOfRestaurantImagePerCommunity = 6

type communityRepository struct {
	db *sqlx.DB
}

func NewCommunityRepository(db *sqlx.DB) models.CommunityRepository {
	return &communityRepository{db}
}

func (c *communityRepository) GetCommunityByID(ctx context.Context, id int64) (*models.Community, error) {
	community, err := models_gen.FindCommunity(ctx, c.db, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to find a community by ID: %w", err)
	}

	numRestaurants, err := community.CommunitiesRestaurants().Count(ctx, c.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to count restaurants related to the community: %w", err)
	}

	userCount, err := community.Affiliations().Count(ctx, c.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to count users related to the community: %w", err)
	}

	urls := make([]string, 0, numOfRestaurantImagePerCommunity)
	err = c.db.SelectContext(ctx, &urls, `SELECT image_url
FROM communities_restaurants
       INNER JOIN restaurants r ON communities_restaurants.restaurant_id = r.id
WHERE community_id = ?
ORDER BY r.created_at DESC
LIMIT ?;`, id, numOfRestaurantImagePerCommunity)
	if err != nil {
		return nil, xerrors.Errorf("failed to get image url from database: %w", err)
	}

	return &models.Community{
		Community:      *community,
		NumRestaurants: int(numRestaurants),
		NumUsers:       int(userCount),
		ImageUrls:      urls,
	}, nil
}

func (c *communityRepository) NewCommunity(ctx context.Context, name string, description string, loc models.Location) (*models.Community, error) {
	result, err := c.db.ExecContext(ctx, "INSERT INTO communities(name, description, latitude,longitude,image_file) VALUES (?,?,?,?,'')",
		name, description, loc.Latitude, loc.Longitude)
	if err != nil {
		return nil, xerrors.Errorf("failed to insert to database: %w", err)
	}

	communityId, err := result.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf("failed to get last insert id: %w", err)
	}
	community, err := models_gen.FindCommunity(ctx, c.db, communityId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get community: %w", err)
	}
	return &models.Community{Community: *community}, nil
}

func (c *communityRepository) SearchCommunity(ctx context.Context, keyword string) ([]*models.Community, error) {
	query := "%" + keyword + "%"
	comm, err := models_gen.Communities(qm.Where("name LIKE ? OR description LIKE ?", query, query)).All(ctx, c.db) // c.db.SelectContext(ctx, &comm, "SELECT * FROM communities WHERE name LIKE ? OR description LIKE ?", query, query)
	if err != nil {
		return nil, xerrors.Errorf("failed to get from database: %w", err)
	}

	if len(comm) == 0 {
		return []*models.Community{}, nil
	}

	commIds := make([]int64, 0, len(comm))
	for _, community := range comm {
		commIds = append(commIds, community.ID)
	}

	sql, params, err := sqlx.In(`SELECT joined.community_id AS community_id, image_url
FROM (
       SELECT community_id,
              r.image_url,
              ROW_NUMBER() OVER (PARTITION BY community_id ORDER BY r.created_at DESC ) AS num
       FROM communities_restaurants
              INNER JOIN restaurants r
                         ON communities_restaurants.restaurant_id = r.id
       ORDER BY community_id, num) AS joined
WHERE joined.num <= ?
  AND community_id IN (?)`, numOfRestaurantImagePerCommunity, commIds)
	if err != nil {
		return nil, xerrors.Errorf("failed to prepare params for image url: %w", err)
	}
	images := make([]communityImage, 0)
	err = c.db.SelectContext(ctx, &images, sql, params...)
	if err != nil {
		return nil, xerrors.Errorf("failed to get image urls from database: %w", err)
	}
	imageUrlMap := toImageUrlMap(images)

	sql, params, err = sqlx.In(`SELECT community_id, COUNT(restaurant_id) AS num
FROM communities_restaurants
WHERE community_id IN (?)
GROUP BY community_id`, commIds)
	if err != nil {
		return nil, xerrors.Errorf("failed to prepare params for image url: %w", err)
	}
	restaurantNums := make([]communityNums, 0, len(comm))
	err = c.db.SelectContext(ctx, &restaurantNums, sql, params...)
	if err != nil {
		return nil, xerrors.Errorf("failed to get image urls from database: %w", err)
	}
	restaurantNumsMap := toCommunityNumsMap(restaurantNums)

	sql, params, err = sqlx.In(`SELECT community_id, COUNT(user_id) AS num
FROM affiliation
WHERE community_id IN (?)
GROUP BY community_id`, commIds)
	if err != nil {
		return nil, xerrors.Errorf("failed to prepare params for image url: %w", err)
	}
	userNums := make([]communityNums, 0, len(comm))
	err = c.db.SelectContext(ctx, &userNums, sql, params...)
	if err != nil {
		return nil, xerrors.Errorf("failed to get image urls from database: %w", err)
	}
	userNumsMap := toCommunityNumsMap(userNums)

	ret := make([]*models.Community, 0)
	for _, community := range comm {
		imgs, ok := imageUrlMap[community.ID]
		if !ok {
			imgs = []string{}
		}
		restaurantNum, ok := restaurantNumsMap[community.ID]
		if !ok {
			restaurantNum = 0
		}
		userNum, ok := userNumsMap[community.ID]
		if !ok {
			userNum = 0
		}

		ret = append(ret, &models.Community{
			Community:      *community,
			ImageUrls:      imgs,
			NumRestaurants: restaurantNum,
			NumUsers:       userNum,
		})
	}
	return ret, nil
}

func (c *communityRepository) UpdateCommunity(ctx context.Context, communityId int64, name string, description string, loc models.Location) (*models.Community, error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return nil, xerrors.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	community, err := models_gen.FindCommunity(ctx, tx, communityId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get community: %w", err)
	}

	community.Name = name
	community.Description = description
	community.Latitude.SetValid(loc.Latitude)
	community.Longitude.SetValid(loc.Longitude)
	count, err := community.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, xerrors.Errorf("failed to update community: %w", err)
	}
	if count != 1 {
		return nil, xerrors.New("# of affected rows is not 1")
	}

	return &models.Community{
		Community: *community,
		ImageUrls: []string{},
	}, nil
}
