package models

import (
	"github.com/jphacks/B_2121_server/openapi"
	"golang.org/x/xerrors"
)

type AuthInfo struct {
	Vendor AuthVendor
	Token  string
}

type AuthVendor string

const (
	AuthVendorGoogle    AuthVendor = "google"
	AuthVendorApple     AuthVendor = "apple"
	AuthVendorAnonymous AuthVendor = "anonymous"
)

func ToOpenApiAuthVendor(v AuthVendor) (openapi.AuthVendor, error) {
	switch v {
	case AuthVendorApple:
		return openapi.AuthVendorApple, nil
	case AuthVendorGoogle:
		return openapi.AuthVendorGoogle, nil
	case AuthVendorAnonymous:
		return openapi.AuthVendorAnonymous, nil
	default:
		return "", xerrors.Errorf("Invalid auth vendor: %s", v)
	}
}

func FromOpenApiAuthVendor(v openapi.AuthVendor) (AuthVendor, error) {
	switch v {
	case openapi.AuthVendorApple:
		return AuthVendorApple, nil
	case openapi.AuthVendorGoogle:
		return AuthVendorGoogle, nil
	case openapi.AuthVendorAnonymous:
		return AuthVendorAnonymous, nil
	default:
		return "", xerrors.Errorf("Invalid auth vendor: %s", v)
	}
}

func (a *AuthInfo) ToOpenApi() (*openapi.AuthInfo, error) {
	vendor, err := ToOpenApiAuthVendor(a.Vendor)
	if err != nil {
		return nil, xerrors.Errorf("failed to convert auth vendor: %w", err)
	}

	return &openapi.AuthInfo{
		Token:  a.Token,
		Vendor: vendor,
	}, nil
}
