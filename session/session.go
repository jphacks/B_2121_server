package session

import (
	"errors"
)

var (
//sessionKey = "session_id"
)

const sessionName = "session_token"

var ErrNoSession = errors.New("no session was found")

type AuthInfo struct {
	Authenticated bool
	UserId        int64
}

//
//type Session struct {
//	*sessions.Session
//}

//func (s *Session) SetUserId(id int32) {
//	s.Values[userIdKey] = id
//}

//func (s *Session) UserId() (int32, error) {
//	v, ok := s.Values[userIdKey]
//	if !ok {
//		return -1, xerrors.New("failed to get user id")
//	}
//
//	id, ok := v.(int32)
//	if !ok {
//		return -1, xerrors.New("failed to get user id")
//	}
//	return id, nil
//}
//
//func (s *Session) SetMaxAge(maxAge int) {
//	s.Options.MaxAge = maxAge
//}

//func NewSession() (string, error) {
//	//cookieStore, err := getSessionStore(ctx)
//	//if err != nil {
//	//	return nil, xerrors.Errorf("failed to fetch session cookieStore: %w", err)
//	//}
//	//r := http.Request{Header: http.Header{}}
//	//s, err := cookieStore.New(&r, sessionKey)
//	//if err != nil {
//	//	return nil, xerrors.Errorf("failed to fetch session cookieStore: %w", err)
//	//}
//	//return &Session{Session: s}, nil
//}

//func Get(token string, store sessions.Store) (*AuthInfo, error) {
//	req := &http.Request{
//		Header: http.Header{},
//	}
//	cookie := &http.Cookie{
//		Name:  sessionName,
//		Value: token,
//	}
//	req.Header.Set("Cookie", cookie.String())
//	sess, err := store.Get(req, sessionName)
//	if err != nil {
//		return nil, err
//	}
//	userId, ok := sess.Values[sessionName]
//	if !ok {
//		return nil, errors.New("user id not found")
//	}
//
//	userIdInt, ok := userId.(int64)
//	if !ok {
//		return nil, errors.New("user id not is invalid")
//	}
//	return &AuthInfo{
//		Authenticated: true,
//		UserId:        userIdInt,
//	}, nil
//}

//func Save(ctx context.Context, session *Session) error {
//	cookieStore, err := getSessionStore(ctx)
//	if err != nil {
//		return xerrors.Errorf("failed to get session cookieStore: %w", err)
//	}
//
//	w := fakeResponseWriter{headers: http.Header{}}
//	err = cookieStore.Save(nil, &w, session.Session)
//	if err != nil {
//		return err
//	}
//
//	cookie := w.headers.Get("set-cookie")
//	if cookie != "" {
//		md := metadata.New(map[string]string{"set-cookie": cookie})
//		err = grpc.SetHeader(ctx, md)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func getCookieText(ctx context.Context) (string, error) {
//	md, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		return "", xerrors.New("failed to get metadata")
//	}
//	cookieTexts := md.Get("cookie")
//
//	// The user agent MUST NOT attach more than one Cookie header field
//	// (Section 5.4 in RFC 6265)
//	if len(cookieTexts) > 1 {
//		return "", xerrors.New("there are two or more cookie headers")
//	}
//
//	if len(cookieTexts) == 0 {
//		return "", http.ErrNoCookie
//	}
//	return cookieTexts[0], nil
//}

//func getSessionStore(ctx context.Context) (sessions.Store, error) {
//	s := ctx.Value(storeKey{})
//	if s == nil {
//		return nil, xerrors.New("session cookieStore not found")
//	}
//	cookieStore, ok := s.(sessions.Store)
//	if !ok {
//		return nil, xerrors.New("failed to get session cookieStore.")
//	}
//	return cookieStore, nil
//}
