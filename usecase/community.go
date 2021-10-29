package usecase

import (
	"context"
	"net/url"

	"github.com/jphacks/B_2121_server/config"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/session"
	"golang.org/x/xerrors"
)

type CommunityUseCase struct {
	sessionStore                   session.Store
	communityRepository            models.CommunityRepository
	affiliationRepository          models.AffiliationRepository
	communityRestaurantsRepository models.CommunityRestaurantsRepository
	imageUrlBase                   string
}

func NewCommunityUseCase(store session.Store, config *config.ServerConfig, communityRepository models.CommunityRepository, affiliationRepository models.AffiliationRepository, communityRestaurantsRepository models.CommunityRestaurantsRepository) CommunityUseCase {
	return CommunityUseCase{
		sessionStore:                   store,
		communityRepository:            communityRepository,
		affiliationRepository:          affiliationRepository,
		communityRestaurantsRepository: communityRestaurantsRepository,
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
