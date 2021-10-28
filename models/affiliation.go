package models

import "context"

type AffiliationRepository interface {
	JoinCommunity(ctx context.Context, userId, communityId int64) error
}
