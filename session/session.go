package session

import (
	"errors"
)

const sessionName = "session_token"

var ErrNoSession = errors.New("no session was found")

type AuthInfo struct {
	Authenticated bool
	UserId        int64
}
