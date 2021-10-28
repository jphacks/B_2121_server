package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

func NewCommentRepository(db *sqlx.DB) models.CommentRepository {
	return &commentRepository{db}
}

type commentRepository struct {
	db *sqlx.DB
}

func (c commentRepository) SetComment(ctx context.Context, communityId, restaurantId int64, comment string) error {
	dbComment, err := models_gen.Comments(qm.Where("community_id = ? AND restaurant_id = ?", communityId, restaurantId)).One(ctx, c.db)
	if xerrors.Is(err, sql.ErrNoRows) {
		result, err := c.db.ExecContext(ctx, "INSERT INTO comments (community_id, restaurant_id, body) VALUES (?, ?, ?)", communityId, restaurantId, comment)
		if err != nil {
			return xerrors.Errorf("failed to get from database: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return xerrors.Errorf("failed to get # of rows affected: %w", err)
		}
		if rows != 1 {
			return xerrors.New("rows affected is not 1")
		}
		return nil
	}
	if err != nil {
		return xerrors.Errorf("failed to get comment from database: %w", err)
	}
	dbComment.Body = comment
	affected, err := dbComment.Update(ctx, c.db, boil.Infer())
	if err != nil {
		return xerrors.Errorf("failed to update database: %w", err)
	}
	if affected != 1 {
		return xerrors.New("# of rows affected is not 1")
	}

	return nil
}

func (c commentRepository) GetComment(ctx context.Context, communityId, restaurantId int64) (string, error) {
	comment, err := models_gen.Comments(qm.Where("community_id = ? AND restaurant_id = ?", communityId, restaurantId)).One(ctx, c.db)
	if xerrors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", xerrors.Errorf("failed to get comment from database: %w", err)
	}

	return comment.Body, err
}
