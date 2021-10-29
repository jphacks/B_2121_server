package usecase

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/gofrs/uuid"
	"github.com/jphacks/B_2121_server/config"
	"github.com/jphacks/B_2121_server/images"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

const profileImageSize = 400

func NewUserUseCase(
	store session.Store,
	userRepo models.UserRepository,
	affiliationRepo models.AffiliationRepository,
	inviteTokenRepository models.InviteTokenRepository,
	config *config.ServerConfig,
) UserUseCase {
	return UserUseCase{
		imageStorePath:        "./profileImages/",
		imageUrlBase:          config.ProfileImageBaseUrl,
		sessionStore:          store,
		userRepo:              userRepo,
		affiliationRepo:       affiliationRepo,
		inviteTokenRepository: inviteTokenRepository,
	}
}

type UserUseCase struct {
	imageStorePath        string
	imageUrlBase          string
	sessionStore          session.Store
	userRepo              models.UserRepository
	affiliationRepo       models.AffiliationRepository
	inviteTokenRepository models.InviteTokenRepository
}

func (u *UserUseCase) UpdateUserProfileImage(ctx context.Context, userId int64, imageData []byte) (imageUrl string, e error) {
	imageId, err := uuid.DefaultGenerator.NewV4()
	if err != nil {
		return "", xerrors.Errorf("failed to generate uuid: %w", err)
	}
	physicalPath := path.Join(u.imageStorePath, imageId.String()+".jpg")
	img, err := images.LoadImage(bytes.NewReader(imageData))
	if err != nil {
		return "", xerrors.Errorf("failed to load image: %w", err)
	}
	img, err = img.CropToSquare()
	if err != nil {
		return "", xerrors.Errorf("failed to crop image: %w", err)
	}
	img = img.ResizeToSquare(profileImageSize)
	file, err := os.Create(physicalPath)
	if err != nil {
		return "", xerrors.Errorf("failed to create image file: %w", err)
	}
	defer func() {
		e1 := file.Close()
		if e1 != nil {
			e = xerrors.Errorf("failed to close file: %w", e1)
		}
	}()
	err = img.Save(file)
	if err != nil {
		return "", xerrors.Errorf("failed to save image: %w", err)
	}
	baseUrl, err := url.Parse(u.imageUrlBase)
	if err != nil {
		return "", xerrors.Errorf("failed to load base url: %w", err)
	}
	baseUrl.Path = path.Join(baseUrl.Path, path.Base(physicalPath))

	err = u.userRepo.UpdateProfileImage(ctx, userId, path.Base(physicalPath))
	if err != nil {
		return "", xerrors.Errorf("failed to update data in db: %w", err)
	}

	return baseUrl.String(), nil
}

func (u *UserUseCase) NewUser(ctx context.Context, name string, authVendor models.AuthVendor) (*models.User, *models.AuthInfo, error) {
	user, err := u.userRepo.NewUser(ctx, name)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to create user to database: %w", err)
	}

	token, err := u.sessionStore.New(user.Id)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to create a new session: %w", err)
	}
	authInfo := &models.AuthInfo{
		Vendor: authVendor,
		Token:  token,
	}
	return user, authInfo, nil
}

func (u *UserUseCase) MyUser(ctx context.Context, userId int64) (*models.UserDetail, error) {
	baseUrl, err := url.Parse(u.imageUrlBase)
	if err != nil {
		return nil, xerrors.Errorf("failed to load base url: %w", err)
	}
	user, err := u.userRepo.GetUserDetailById(ctx, userId, *baseUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to get user detail: %w", err)
	}
	return user, nil
}

func (u *UserUseCase) ListUserCommunities(ctx context.Context, userId int64) ([]*models.Community, error) {
	comm, err := u.userRepo.ListUserCommunity(ctx, userId)
	if err != nil {
		return nil, xerrors.Errorf("failed to list communities: %w", err)
	}
	return comm, nil
}

func (u *UserUseCase) JoinCommunity(ctx context.Context, userId int64, token string) error {
	communityId, err := u.inviteTokenRepository.Verify(ctx, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	return u.affiliationRepo.JoinCommunity(ctx, userId, communityId)
}

func (u *UserUseCase) LeaveCommunity(ctx context.Context, userId int64, communityId int64) error {
	return u.affiliationRepo.LeaveCommunity(ctx, userId, communityId)
}

func (u *UserUseCase) UpdateMe(ctx context.Context, userId int64, userName *string) (*models.User, error) {
	baseUrl, err := url.Parse(u.imageUrlBase)
	if err != nil {
		return nil, xerrors.Errorf("failed to load base url: %w", err)
	}
	user, err := u.userRepo.UpdateUser(ctx, &models.UpdateUserInput{
		Id:   userId,
		Name: userName,
	}, *baseUrl)
	if err != nil {
		return nil, err
	}

	return user, nil
}
