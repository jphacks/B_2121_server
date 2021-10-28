package models

import (
	"context"
	"net/url"
)

type AffiliationRepository interface {
	JoinCommunity(ctx context.Context, userId, communityId int64) error
	ListCommunityUsers(ctx context.Context, communityId int64, profileBaseUrl *url.URL) ([]*User, error)
}
