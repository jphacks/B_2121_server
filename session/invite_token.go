package session

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jphacks/B_2121_server/models"
	"golang.org/x/xerrors"
)

const (
	inviteTokenExpiresIn   = 24 * time.Hour
	inviteTokenSessionKey  = "invite_token"
	inviteTokenSessionName = "invite_token"
)

func NewInviteTokenRepository() models.InviteTokenRepository {
	cookieStore := sessions.NewCookieStore([]byte(inviteTokenSessionKey))
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.SameSite = http.SameSiteLaxMode
	cookieStore.Options.MaxAge = int(inviteTokenExpiresIn.Seconds())
	cookieStore.Options.Secure = false
	return &inviteTokenRepository{
		cookieStore: cookieStore,
	}
}

type inviteTokenRepository struct {
	cookieStore sessions.Store
}

func (r *inviteTokenRepository) Issue(ctx context.Context, communityId int64) (*models.InviteToken, error) {
	req := http.Request{Header: http.Header{}}
	newSession, err := r.cookieStore.New(&req, inviteTokenSessionName)
	if err != nil {
		return nil, xerrors.Errorf("failed to fetch the store for invite tokens: %w", err)
	}
	newSession.Values[inviteTokenSessionName] = communityId
	w := fakeResponseWriter{headers: http.Header{}}
	err = r.cookieStore.Save(nil, &w, newSession)
	if err != nil {
		return nil, err
	}

	cookie := w.headers.Get("set-cookie")
	token := getCookieValue(cookie)
	return &models.InviteToken{
		Token:     token,
		ExpiresIn: inviteTokenExpiresIn,
	}, nil
}

func (r *inviteTokenRepository) Verify(ctx context.Context, token string) (string, error) {
	panic("implement me")
}
