package db

import (
	"context"
	"database/sql"
	"net/url"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/jphacks/B_2121_server/models"
	"github.com/jphacks/B_2121_server/models_gen"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/xerrors"
)

type affiliationRepository struct {
	db *sqlx.DB
}

func NewAffiliationRepository(db *sqlx.DB) models.AffiliationRepository {
	return &affiliationRepository{db}
}

func (a affiliationRepository) JoinCommunity(ctx context.Context, userId, communityId int64) error {
	result, err := a.db.ExecContext(ctx, "INSERT INTO affiliation(community_id, user_id) VALUES (?, ?)", communityId, userId)
	if err != nil {
		return xerrors.Errorf("failed to add affiliation: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return xerrors.Errorf("failed to egt row affected: %w", err)
	}

	if rowsAffected != 1 {
		return xerrors.New("rows affected is not 1")
	}

	return nil
}

func (a affiliationRepository) ListCommunityUsers(ctx context.Context, communityId int64, profileBaseUrl *url.URL) ([]*models.User, error) {
	users, err := models_gen.Users(
		qm.InnerJoin("affiliation ON affiliation.user_id = users.id"),
		qm.Where("community_id = ?", communityId),
	).All(ctx, a.db)
	if err != nil {
		return nil, xerrors.Errorf("failed to get from database: %w", err)
	}

	ret := make([]*models.User, 0)
	for _, user := range users {
		if user.ProfileImageFile.Valid {
			imageUrlBase := *profileBaseUrl
			imageUrlBase.Path = path.Join(imageUrlBase.Path, user.ProfileImageFile.String)
			ret = append(ret, &models.User{
				Id:              user.ID,
				Name:            user.Name,
				ProfileImageUrl: imageUrlBase.String(),
			})
		} else {
			ret = append(ret, &models.User{
				Id:              user.ID,
				Name:            user.Name,
				ProfileImageUrl: "",
			})
		}
	}

	return ret, nil
}

func (a affiliationRepository) LeaveCommunity(ctx context.Context, userId, communityId int64) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	aff, err := models_gen.Affiliations(qm.Where("user_id = ?", userId), qm.And("community_id = ?", communityId)).One(ctx, tx)
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return xerrors.Errorf("affiliation not found: %w", err)
		}
		return err
	}

	rowsAffected, err := aff.Delete(ctx, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return xerrors.Errorf("failed to get row affected: %w", err)
	}
	if rowsAffected != 1 {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return xerrors.New("rows affected is not 1")
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
