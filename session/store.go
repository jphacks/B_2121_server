package session

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/xerrors"
)

const cookieMaxAge = 3600 * 24 * 365 * 10 // 10 years

type Store interface {
	Get(token string) (*AuthInfo, error)
	New(userId int64) (string, error)
}

type store struct {
	cookieStore sessions.Store
}

func (s store) Get(token string) (*AuthInfo, error) {
	req := &http.Request{
		Header: http.Header{},
	}
	cookie := &http.Cookie{
		Name:  sessionName,
		Value: token,
	}
	req.Header.Set("Cookie", cookie.String())
	sess, err := s.cookieStore.Get(req, sessionName)
	if err != nil {
		return nil, err
	}
	userId, ok := sess.Values[sessionName]
	if !ok {
		return nil, errors.New("user id not found")
	}

	userIdInt, ok := userId.(int64)
	if !ok {
		return nil, errors.New("user id not is invalid")
	}
	return &AuthInfo{
		Authenticated: true,
		UserId:        userIdInt,
	}, nil
}

func (s store) New(userId int64) (string, error) {
	r := http.Request{Header: http.Header{}}
	newSession, err := s.cookieStore.New(&r, sessionName)
	if err != nil {
		return "", xerrors.Errorf("failed to fetch session cookieStore: %w", err)
	}
	newSession.Values[sessionName] = userId

	w := fakeResponseWriter{headers: http.Header{}}
	err = s.cookieStore.Save(nil, &w, newSession)
	if err != nil {
		return "", err
	}

	cookie := w.headers.Get("set-cookie")
	return getCookieValue(cookie), nil
}

func NewStore(key string) Store {
	cookieStore := sessions.NewCookieStore([]byte(key))
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.SameSite = http.SameSiteLaxMode
	cookieStore.Options.MaxAge = cookieMaxAge
	cookieStore.Options.Secure = false
	return &store{cookieStore: cookieStore}
}

func getCookieValue(cookieText string) string {
	header := http.Header{}
	header.Add("Cookie", cookieText)
	request := http.Request{Header: header}
	return request.Cookies()[0].Value
}
