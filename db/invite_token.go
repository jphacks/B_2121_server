package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/lithammer/shortuuid/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"
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
	tokenDigest, err := newDigestString(token)
	if err != nil {
		return nil, err
	}

	var it models_gen.InviteToken
	it.TokenDigest = tokenDigest
	it.CommunityID = communityId
	it.ExpiresAt = time.Now().Add(inviteTokenExpiresIn)
	err = it.Insert(ctx, r.db, boil.Whitelist("token_digest", "community_id", "expires_at"))
	if err != nil {
		return nil, err
	}

	return &models.InviteToken{
		Token:     token,
		ExpiresIn: inviteTokenExpiresIn,
	}, nil
}

func (r *inviteTokenRepository) Verify(ctx context.Context, token string) (int64, error) {
	tokenDigest, err := newDigestString(token)
	if err != nil {
		return 0, err
	}
	it, err := models_gen.FindInviteToken(ctx, r.db, tokenDigest, "community_id", "expires_at")
	if err != nil {
		return 0, err
	}

	if it.ExpiresAt.Before(time.Now()) {
		return 0, xerrors.New("The token already expired")
	}

	return it.CommunityID, nil
}

func newDigestString(token string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(token), 10)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
