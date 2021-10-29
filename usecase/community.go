package usecase

import (
	"context"
	"net/http"
	"net/url"

	"github.com/jphacks/B_2121_server/config"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

type CommunityUseCase struct {
	sessionStore                   session.Store
	communityRepository            models.CommunityRepository
	affiliationRepository          models.AffiliationRepository
	communityRestaurantsRepository models.CommunityRestaurantsRepository
	inviteTokenRepository          models.InviteTokenRepository
	userRepository                 models.UserRepository
	imageUrlBase                   string
}

func NewCommunityUseCase(
	store session.Store,
	config *config.ServerConfig,
	communityRepository models.CommunityRepository,
	affiliationRepository models.AffiliationRepository,
	communityRestaurantsRepository models.CommunityRestaurantsRepository,
	inviteTokenRepository models.InviteTokenRepository,
	userRepository models.UserRepository,
) CommunityUseCase {
	return CommunityUseCase{
		sessionStore:                   store,
		communityRepository:            communityRepository,
		affiliationRepository:          affiliationRepository,
		communityRestaurantsRepository: communityRestaurantsRepository,
		inviteTokenRepository:          inviteTokenRepository,
		userRepository:                 userRepository,
		imageUrlBase:                   config.ProfileImageBaseUrl,
	}
}

func (u *CommunityUseCase) GetCommunity(ctx context.Context, id int64) (*models.Community, error) {
	return u.communityRepository.GetCommunityByID(ctx, id)
}

func (u *CommunityUseCase) NewCommunity(ctx context.Context, userId int64, name string, description string, loc models.Location) (*models.Community, error) {
	community, err := u.communityRepository.NewCommunity(ctx, name, description, loc)
	if err != nil {
		return nil, xerrors.Errorf("failed to create community: %w", err)
	}

	err = u.affiliationRepository.JoinCommunity(ctx, userId, community.ID)
	if err != nil {
		return nil, xerrors.Errorf("failed to join community: %w", err)
	}
	community.NumUsers = 1
	community.ImageUrls = []string{}
	return community, nil
}

func (u *CommunityUseCase) SearchCommunity(ctx context.Context, keyword string) ([]*models.Community, error) {
	comm, err := u.communityRepository.SearchCommunity(ctx, keyword)
	if err != nil {
		return nil, xerrors.Errorf("failed to search community: %w", err)
	}

	return comm, nil
}

func (u *CommunityUseCase) ListRestaurants(ctx context.Context, communityId int64) ([]*models.Restaurant, error) {
	rest, err := u.communityRestaurantsRepository.ListCommunityRestaurants(ctx, communityId)
	if err != nil {
		return nil, xerrors.Errorf("failed to get restaurants of community: %w", err)
	}
	return rest, nil
}

func (u *CommunityUseCase) ListUsers(ctx context.Context, communityId int64) ([]*models.User, error) {
	baseUrl, err := url.Parse(u.imageUrlBase)
	if err != nil {
		return nil, xerrors.Errorf("failed to load base url: %w", err)
	}

	users, err := u.affiliationRepository.ListCommunityUsers(ctx, communityId, baseUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

func (u *CommunityUseCase) IssueInviteToken(ctx context.Context, issuerId int64, communityId int64) (*models.InviteToken, error) {
	permitted, err := u.userRepository.ExistInCommunity(ctx, issuerId, communityId)
	if err != nil {
		return nil, err
	}

	if !permitted {
		return nil, echo.NewHTTPError(http.StatusForbidden)
	}

	inviteToken, err := u.inviteTokenRepository.Issue(ctx, communityId)
	if err != nil {
		return nil, err
	}

	return inviteToken, nil
}
