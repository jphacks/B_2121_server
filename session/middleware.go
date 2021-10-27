package session

import (
	"github.com/labstack/echo/v4"
)

const authInfoKey = "authInfo"

type MiddlewareConfig struct {
	SessionStore Store
}

func GetAuthInfo(ctx echo.Context) *AuthInfo {
	info := ctx.Get(authInfoKey)
	if info == nil {
		return &AuthInfo{
			Authenticated: false,
		}
	}
	authInfo, ok := info.(AuthInfo)
	if !ok {
		return &AuthInfo{
			Authenticated: false,
		}
	}
	return &authInfo
}

func NewSessionMiddleware(config *MiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			auth := context.Request().Header.Get(echo.HeaderAuthorization)
			if len(auth) == 0 {
				context.Set(authInfoKey, AuthInfo{
					Authenticated: false,
					UserId:        0,
				})
				return next(context)
			}
			info, err := config.SessionStore.Get(auth)
			if err != nil {
				info = &AuthInfo{
					Authenticated: false,
					UserId:        0,
				}
			}
			context.Set(authInfoKey, *info)
			return next(context)
		}
	}
}
