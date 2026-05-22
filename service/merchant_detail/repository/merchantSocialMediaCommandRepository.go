package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_social_link_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type merchantSocialLinkCommandRepository struct {
	pool *pgxpool.Pool
	db   *db.Queries
}

func NewMerchantSocialLinkCommandRepository(pool *pgxpool.Pool, db *db.Queries) *merchantSocialLinkCommandRepository {
	return &merchantSocialLinkCommandRepository{
		db:   db,
		pool: pool,
	}
}

func (r *merchantSocialLinkCommandRepository) CreateSocialLink(ctx context.Context, req *requests.CreateBatchMerchantSocialRequest) ([]*db.CreateMerchantSocialMediaLinkRow, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, merchant_social_link_errors.ErrBeginTx
	}

	defer tx.Rollback(ctx)

	qtx := r.db.WithTx(tx)

	results := make([]*db.CreateMerchantSocialMediaLinkRow, 0, len(req.SocialLinks))

	for _, link := range req.SocialLinks {
		params := db.CreateMerchantSocialMediaLinkParams{
			MerchantDetailID: int32(req.MerchantDetailID),
			Platform:         link.Platform,
			Url:              link.Url,
		}

		res, err := qtx.CreateMerchantSocialMediaLink(ctx, params)
		if err != nil {
			return nil, merchant_social_link_errors.ErrCreateMerchantSocialLink
		}

		results = append(results, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, merchant_social_link_errors.ErrCommitTx
	}

	return results, nil
}

func (r *merchantSocialLinkCommandRepository) UpdateSocialLink(ctx context.Context, req *requests.UpdateBatchMerchantSocialRequest) ([]*db.UpdateMerchantSocialMediaLinkRow, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, merchant_social_link_errors.ErrBeginTx
	}
	defer tx.Rollback(ctx)

	qtx := r.db.WithTx(tx)

	results := make([]*db.UpdateMerchantSocialMediaLinkRow, 0, len(req.SocialLinks))

	for _, link := range req.SocialLinks {
		params := db.UpdateMerchantSocialMediaLinkParams{
			MerchantSocialID: int32(link.ID),
			Platform:         link.Platform,
			Url:              link.Url,
		}

		res, err := qtx.UpdateMerchantSocialMediaLink(ctx, params)
		if err != nil {
			return nil, merchant_social_link_errors.ErrUpdateMerchantSocialLink
		}

		results = append(results, res)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, merchant_social_link_errors.ErrCommitTx
	}

	return results, nil
}

func (r *merchantSocialLinkCommandRepository) TrashSocialLink(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.TrashMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchant_social_link_errors.ErrTrashMerchantSocialLink.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreSocialLink(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.RestoreMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchant_social_link_errors.ErrRestoreMerchantSocialLink.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeletePermanentSocialLink(ctx context.Context, socialID int) (bool, error) {
	err := r.db.DeleteMerchantSocialMediaLinkPermanently(ctx, int32(socialID))
	if err != nil {
		return false, merchant_social_link_errors.ErrDeletePermanentMerchantSocialLink.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreAllSocialLinks(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantSocialMediaLinks(ctx)
	if err != nil {
		return false, merchant_social_link_errors.ErrRestoreAllMerchantSocialLinks.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllMerchantSocialMediaLinksPermanently(ctx)
	if err != nil {
		return false, merchant_social_link_errors.ErrDeleteAllPermanentMerchantSocialLinks.WithInternal(err)
	}

	return true, nil
}
