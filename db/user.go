package db

import (
	"context"
	"database/sql"
	"net/url"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) models.UserRepository {
	return &userRepository{db: db}
}

func (u userRepository) GetUserById(ctx context.Context, id int64, profileImageBase url.URL) (*models.User, error) {
	user, err := models_gen.FindUser(ctx, u.db, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to find a user by ID: %w", err)
	}

	return fromGenUser(user, profileImageBase), err
}

func (u userRepository) GetUserDetailById(ctx context.Context, id int64, profileImageBase url.URL) (*models.UserDetail, error) {
	user, err := models_gen.FindUser(ctx, u.db, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to find a user by ID: %w", err)
	}

	communityCount, err := user.Affiliations().Count(ctx, u.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to count joined communities: %w", err)
	}
	bookmarkCount, err := user.Bookmarks().Count(ctx, u.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to count bookmarked communities: %w", err)
	}

	modelUser := fromGenUser(user, profileImageBase)
	return &models.UserDetail{
		User:           *modelUser,
		BookmarkCount:  int(bookmarkCount),
		CommunityCount: int(communityCount),
	}, err
}

func (u userRepository) NewUser(ctx context.Context, userName string) (*models.User, error) {
	result, err := u.db.ExecContext(ctx, "INSERT INTO users(`name`) VALUES (?)", userName)
	if err != nil {
		return nil, xerrors.Errorf("failed to insert a new user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf("failed to get last insert id: %w", err)
	}
	return &models.User{
		Id:   id,
		Name: userName,
	}, nil
}

func (u userRepository) UpdateProfileImage(ctx context.Context, userId int64, fileName string) error {
	user, err := models_gen.FindUser(ctx, u.db, userId)
	if err != nil {
		return xerrors.Errorf("failed to find a user by ID: %w", err)
	}
	user.ProfileImageFile.SetValid(fileName)
	_, err = user.Update(ctx, u.db, boil.Infer())
	if err != nil {
		return xerrors.Errorf("failed to update database: %w", err)
	}
	return nil
}

type communityImage struct {
	CommunityId int64  `db:"community_id"`
	ImageUrl    string `db:"image_url"`
}

type communityNums struct {
	CommunityId int64 `db:"community_id"`
	Num         int   `db:"num"`
}

func (u userRepository) ListUserCommunity(ctx context.Context, userId int64) ([]*models.Community, error) {
	community, err := models_gen.Communities(
		qm.InnerJoin("affiliation ON affiliation.community_id = communities.id"),
		qm.Where("user_id=?", userId),
		qm.OrderBy("id DESC"),
	).All(ctx, u.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to get communities: %w", err)
	}

	images := make([]communityImage, 0)
	err = u.db.SelectContext(ctx, &images, `SELECT joined.community_id AS community_id, image_url
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
       INNER JOIN affiliation ON joined.community_id = affiliation.community_id
WHERE joined.num <= ?
  AND user_id = ?`, numOfRestaurantImagePerCommunity, userId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get image urls from database: %w", err)
	}
	imageUrlMap := toImageUrlMap(images)

	restaurantNums := make([]communityNums, 0, len(community))
	err = u.db.SelectContext(ctx, &restaurantNums, `SELECT a.community_id AS community_id, COUNT(restaurant_id) AS num
FROM communities_restaurants INNER JOIN affiliation a ON communities_restaurants.community_id = a.community_id
WHERE user_id = ?
GROUP BY community_id;`, userId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get restaurant count from database: %s", err)
	}
	restaurantNumsMap := toCommunityNumsMap(restaurantNums)

	userNums := make([]communityNums, 0, len(community))
	err = u.db.SelectContext(ctx, &userNums, `SELECT a1.community_id AS community_id, COUNT(a1.user_id) AS num
FROM affiliation AS a1 INNER JOIN affiliation AS a2 ON a1.community_id = a2.community_id
WHERE a2.user_id = ?
GROUP BY a1.community_id`, userId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get user count from database: %s", err)
	}
	userNumsMap := toCommunityNumsMap(userNums)

	ret := make([]*models.Community, 0)
	for _, c := range community {
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

func (u userRepository) ExistInCommunity(ctx context.Context, userId int64, communityId int64) (bool, error) {
	num, err := models_gen.Affiliations(qm.Where("community_id = ? AND user_id = ?", communityId, userId)).Count(ctx, u.db)
	if err != nil {
		return false, xerrors.Errorf("failed to get affiliations from database: %w", err)
	}

	return num > 0, nil
}

func (u userRepository) UpdateUser(ctx context.Context, input *models.UpdateUserInput, profileImageBase url.URL) (*models.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	user, err := models_gen.FindUser(ctx, tx, input.Id)
	if err != nil {
		return nil, err
	}

	updatedColumns := make([]string, 0)

	if input.Name != nil {
		updatedColumns = append(updatedColumns, "name")
		user.Name = *input.Name
	}

	if len(updatedColumns) == 0 {
		return fromGenUser(user, profileImageBase), nil
	}

	rows, err := user.Update(ctx, tx, boil.Whitelist(updatedColumns...))
	if err != nil {
		return nil, err
	}

	if rows != 1 {
		return nil, xerrors.New("rows affected is not 1")
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return fromGenUser(user, profileImageBase), nil
}

func fromGenUser(u *models_gen.User, imageUrlBase url.URL) *models.User {
	if u.ProfileImageFile.Valid {
		imageUrlBase.Path = path.Join(imageUrlBase.Path, u.ProfileImageFile.String)
		return &models.User{
			Id:              u.ID,
			Name:            u.Name,
			ProfileImageUrl: imageUrlBase.String(),
		}
	}

	return &models.User{
		Id:              u.ID,
		Name:            u.Name,
		ProfileImageUrl: "",
	}
}

func toImageUrlMap(imageUrls []communityImage) map[int64][]string {
	ret := make(map[int64][]string)
	for _, imageUrl := range imageUrls {
		arr, ok := ret[imageUrl.CommunityId]
		if ok {
			ret[imageUrl.CommunityId] = append(arr, imageUrl.ImageUrl)
		} else {
			ret[imageUrl.CommunityId] = []string{imageUrl.ImageUrl}
		}
	}
	return ret
}

func toCommunityNumsMap(nums []communityNums) map[int64]int {
	ret := make(map[int64]int)
	for _, imageUrl := range nums {
		ret[imageUrl.CommunityId] = imageUrl.Num
	}
	return ret
}
