package usecase

import (
	"context"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/session"
	"golang.org/x/xerrors"
)

type CommunityUseCase struct {
	sessionStore          session.Store
	communityRepository   models.CommunityRepository
	affiliationRepository models.AffiliationRepository
}

func NewCommunityUseCase(store session.Store, communityRepository models.CommunityRepository, affiliationRepository models.AffiliationRepository) CommunityUseCase {
	return CommunityUseCase{sessionStore: store, communityRepository: communityRepository, affiliationRepository: affiliationRepository}
}

func (u *CommunityUseCase) GetCommunity(ctx context.Context, id int64) (*models.CommunityDetail, error) {
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
	return community, nil
}

func (u *CommunityUseCase) SearchCommunity(ctx context.Context, keyword string) ([]*models.Community, error) {
	comm, err := u.communityRepository.SearchCommunity(ctx, keyword)
	if err != nil {
		return nil, xerrors.Errorf("failed to search community: %w", err)
	}

	return comm, nil
}
