package models

import (
	"context"
	"net/url"

	"github.com/jphacks/B_2121_server/openapi"
)

type User struct {
	Id              int64
	Name            string
	ProfileImageUrl string
}

func (u *User) ToOpenApiUser() *openapi.User {
	return &openapi.User{
		Id:              openapi.Long(u.Id),
		Name:            u.Name,
		ProfileImageUrl: &u.ProfileImageUrl,
	}
}

type UserDetail struct {
	User
	BookmarkCount  int
	CommunityCount int
}

func (u *UserDetail) ToOpenApi() *openapi.UserDetail {
	return &openapi.UserDetail{
		User:           *u.User.ToOpenApiUser(),
		BookmarkCount:  u.BookmarkCount,
		CommunityCount: u.CommunityCount,
	}
}

type UpdateUserInput struct {
	Id   int64
	Name *string
}

type UserRepository interface {
	GetUserById(ctx context.Context, id int64, profileImageBase url.URL) (*User, error)
	NewUser(ctx context.Context, userName string) (*User, error)
	GetUserDetailById(ctx context.Context, id int64, profileImageBase url.URL) (*UserDetail, error)
	UpdateProfileImage(ctx context.Context, userId int64, fileName string) error
	ListUserCommunity(ctx context.Context, userId int64) ([]*Community, error)
	ExistInCommunity(ctx context.Context, userId int64, communityId int64) (bool, error)
	UpdateUser(ctx context.Context, input *UpdateUserInput, profileImageBase url.URL) (*User, error)
}
