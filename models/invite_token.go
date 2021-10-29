package models

import (
	"context"
	"time"
)

type InviteToken struct {
	Token     string
	ExpiresIn time.Duration
}

type InviteTokenRepository interface {
	Issue(ctx context.Context, communityId int64) (*InviteToken, error)
	// Verify returns communityId if the token verified
	Verify(ctx context.Context, token string) (string, error)
}
