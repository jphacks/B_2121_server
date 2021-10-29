package db

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/lithammer/shortuuid/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/xerrors"
)

const (
	inviteTokenExpiresIn = 24 * time.Hour
)

type inviteTokenRepository struct {
	db *sqlx.DB
}

func NewInviteTokenRepository(db *sqlx.DB) models.InviteTokenRepository {
	return &inviteTokenRepository{db: db}
}

func (r *inviteTokenRepository) Issue(ctx context.Context, communityId int64) (*models.InviteToken, error) {
	token := shortuuid.New()

	var it models_gen.InviteToken
	it.TokenDigest = newDigestString(token)
	it.CommunityID = communityId
	it.ExpiresAt = time.Now().Add(inviteTokenExpiresIn)
	err := it.Insert(ctx, r.db, boil.Whitelist("token_digest", "community_id", "expires_at"))
	if err != nil {
		return nil, err
	}

	return &models.InviteToken{
		Token:     token,
		ExpiresIn: inviteTokenExpiresIn,
	}, nil
}

func (r *inviteTokenRepository) Verify(ctx context.Context, token string) (int64, error) {
	it, err := models_gen.FindInviteToken(ctx, r.db, newDigestString(token), "community_id", "expires_at")
	if err != nil {
		return 0, err
	}

	if it.ExpiresAt.Before(time.Now()) {
		return 0, xerrors.New("The token already expired")
	}

	return it.CommunityID, nil
}

func newDigestString(token string) string {
	b := sha256.Sum256([]byte(token))
	return string(b[:])
}
