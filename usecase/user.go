package usecase

import (
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/session"
	"golang.org/x/xerrors"
)

type UserUseCase interface {
	NewUser(name string, authVendor models.AuthVendor) (*models.User, *models.AuthInfo, error)
	MyUser(userId int64) (*models.UserDetail, error)
}

func NewUserUseCase(store session.Store) UserUseCase {
	return &userUseCase{
		sessionStore: store,
	}
}

type userUseCase struct {
	sessionStore session.Store
}

func (u *userUseCase) NewUser(name string, authVendor models.AuthVendor) (*models.User, *models.AuthInfo, error) {
	// TODO: update id, save user info
	var userId int64 = 1
	user := &models.User{
		Id:              userId,
		Name:            name,
		ProfileImageUrl: "",
	}
	token, err := u.sessionStore.New(userId)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to create a new session: %w", err)
	}
	authInfo := &models.AuthInfo{
		Vendor: authVendor,
		Token:  token,
	}
	return user, authInfo, nil
}

func (u *userUseCase) MyUser(userId int64) (*models.UserDetail, error) {
	// TODO: This is a test implementation !!!
	return &models.UserDetail{
		User: models.User{
			Name: "your name",
			Id:   userId,
		},
		BookmarkCount:  0,
		CommunityCount: 0,
	}, nil
}
