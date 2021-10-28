package usecase

import (
	"context"

	"github.com/jphacks/B_2121_server/models"
)

func NewBookmarkUseCase(bookmarkRepo models.BookmarkRepository) BookmarkUseCase {
	return BookmarkUseCase{bookmarkRepo: bookmarkRepo}
}

type BookmarkUseCase struct {
	bookmarkRepo models.BookmarkRepository
}

func (u *BookmarkUseCase) CreateBookmark(ctx context.Context, userId, communityId int64) error {
	_, err := u.bookmarkRepo.CreateBookmark(ctx, userId, communityId)
	if err != nil {
		return err
	}
	return nil
}
