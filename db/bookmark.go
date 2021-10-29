package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

type bookmarkRepository struct {
	db *sqlx.DB
}

func (b *bookmarkRepository) ListBookmarkByUserId(ctx context.Context, userId int64) ([]models.Community, error) {
	comms, err := models_gen.Communities(
		qm.InnerJoin("bookmarks ON bookmarks.community_id = communities.id"),
		qm.Where("user_id = ?", userId),
	).All(ctx, b.db)
	if err != nil {
		return nil, err
	}

	images := make([]communityImage, 0)
	err = b.db.SelectContext(ctx, &images, `SELECT joined.community_id, image_url
FROM (
       SELECT communities_restaurants.id,
              community_id,
              restaurant_id,
              r.image_url,
              communities_restaurants.created_at,
              ROW_NUMBER() OVER (PARTITION BY community_id ORDER BY r.created_at DESC ) AS num
       FROM communities_restaurants
              INNER JOIN restaurants r
                         ON communities_restaurants.restaurant_id = r.id
       ORDER BY community_id, num) AS joined
       INNER JOIN bookmarks ON joined.community_id = bookmarks.community_id
WHERE joined.num <= ?
  AND user_id = ?`, numOfRestaurantImagePerCommunity, userId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get image urls from database: %w", err)
	}
	imageUrlMap := toImageUrlMap(images)

	restaurantNums := make([]communityNums, 0, len(comms))
	err = b.db.SelectContext(ctx, &restaurantNums, `SELECT b.community_id as community_id, count(restaurant_id) as num
FROM communities_restaurants INNER JOIN bookmarks b ON communities_restaurants.community_id = b.community_id
WHERE user_id = ?
GROUP BY community_id`, userId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get restaurant count from database: %s", err)
	}
	restaurantNumsMap := toCommunityNumsMap(restaurantNums)

	userNums := make([]communityNums, 0, len(comms))
	err = b.db.SelectContext(ctx, &userNums, `SELECT a1.community_id as community_id, count(a1.user_id) as num
FROM affiliation as a1 INNER JOIN bookmarks as b ON a1.community_id = b.community_id
WHERE b.user_id = ?
GROUP BY a1.community_id`, userId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get user count from database: %s", err)
	}
	userNumsMap := toCommunityNumsMap(userNums)

	ret := make([]models.Community, 0, len(comms))
	for _, comm := range comms {
		images, ok := imageUrlMap[comm.ID]
		if !ok {
			images = []string{}
		}
		restaurantNum, ok := restaurantNumsMap[comm.ID]
		if !ok {
			restaurantNum = 0
		}

		userNum, ok := userNumsMap[comm.ID]
		if !ok {
			userNum = 0
		}

		ret = append(ret, models.Community{
			Community:      *comm,
			ImageUrls:      images,
			NumRestaurants: restaurantNum,
			NumUsers:       userNum,
		})
	}

	return ret, nil
}

func (b *bookmarkRepository) CreateBookmark(ctx context.Context, userId, communityId int64) error {
	result, err := b.db.ExecContext(ctx, `INSERT INTO bookmarks(community_id, user_id) VALUES (?, ?)`, communityId, userId)
	if err != nil {
		return xerrors.Errorf("failed to create a bookmark: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return xerrors.Errorf("failed to egt row affected: %w", err)
	}

	if rowsAffected != 1 {
		return xerrors.New("rows affected is not 1")
	}

	return nil
}

func (b *bookmarkRepository) DeleteBookmark(ctx context.Context, userId int64, communityId int64) error {
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	bm, err := models_gen.Bookmarks(qm.Where("user_id = ?", userId), qm.And("community_id = ?", communityId)).One(ctx, tx)
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return xerrors.Errorf("affiliation not found: %w", err)
		}
		return err
	}

	rowsAffected, err := bm.Delete(ctx, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return xerrors.Errorf("failed to get row affected: %w", err)
	}
	if rowsAffected != 1 {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return xerrors.New("rows affected is not 1")
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func NewBookmarkRepository(db *sqlx.DB) models.BookmarkRepository {
	return &bookmarkRepository{db}
}
