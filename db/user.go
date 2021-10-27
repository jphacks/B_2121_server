package db

import (
	"context"
	"net/url"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
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
