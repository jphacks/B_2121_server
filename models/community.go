package models

import (
	"context"

	"github.com/jphacks/B_2121_server/models_gen"
)

type CommunityRepository interface {
	GetCommunityByID(ctx context.Context, id int64) (*models_gen.Community, error)
}
