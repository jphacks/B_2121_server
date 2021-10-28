package usecase

import (
	"context"

	"github.com/jphacks/B_2121_server/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

func NewCommentUseCase(commentRepository models.CommentRepository, userRepository models.UserRepository) CommentUseCase {
	return CommentUseCase{
		commentRepository: commentRepository,
		userRepository:    userRepository,
	}
}

type CommentUseCase struct {
	commentRepository models.CommentRepository
	userRepository    models.UserRepository
}

func (c *CommentUseCase) SetComment(ctx context.Context, userId, communityId, restaurantId int64, comment string) error {
	exits, err := c.userRepository.ExistInCommunity(ctx, userId, communityId)
	if err != nil {
		return xerrors.Errorf("failed to get user existence: %w", err)
	}

	if !exits {
		return echo.ErrForbidden
	}

	err = c.commentRepository.SetComment(ctx, communityId, restaurantId, comment)
	if err != nil {
		return xerrors.Errorf("failed to update comment: %w", err)
	}
	return nil
}

func (c *CommentUseCase) GetComment(ctx context.Context, userId, communityId, restaurantId int64) (string, error) {
	exits, err := c.userRepository.ExistInCommunity(ctx, userId, communityId)
	if err != nil {
		return "", xerrors.Errorf("failed to get user existence: %w", err)
	}

	if !exits {
		return "", echo.ErrForbidden
	}

	comm, err := c.commentRepository.GetComment(ctx, communityId, restaurantId)
	if err != nil {
		return "", xerrors.Errorf("failed to update comment: %w", err)
	}
	return comm, nil
}
