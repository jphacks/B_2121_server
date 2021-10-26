package models

import "github.com/jphacks/B_2121_server/openapi"

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
