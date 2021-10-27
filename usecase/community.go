package usecase

import (
	"context"

	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/session"
)

type CommunityUseCase struct {
	sessionStore        session.Store
	communityRepository models.CommunityRepository
}

func NewCommunityUseCase(store session.Store, communityRepository models.CommunityRepository) CommunityUseCase {
	return CommunityUseCase{sessionStore: store, communityRepository: communityRepository}
}

func (u *CommunityUseCase) GetCommunity(ctx context.Context, id int64) (*models.CommunityDetail, error) {
	return u.communityRepository.GetCommunityByID(ctx, id)
}
