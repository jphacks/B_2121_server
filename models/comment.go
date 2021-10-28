package models

import "context"

type CommentRepository interface {
	SetComment(ctx context.Context, communityId, restaurantId int64, comment string) error
	GetComment(ctx context.Context, communityId, restaurantId int64) (string, error)
}
