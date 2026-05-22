package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type MerchantDetailQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantDetailRow, error)
	FindByIDTrashed(ctx context.Context, user_id int) (*db.GetMerchantDetailTrashedRow, error)
}

type MerchantDetailCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*db.CreateMerchantDetailRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*db.UpdateMerchantDetailRow, error)
	Trash(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	Restore(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandRepository interface {
	CreateSocialLink(ctx context.Context, req *requests.CreateBatchMerchantSocialRequest) ([]*db.CreateMerchantSocialMediaLinkRow, error)
	UpdateSocialLink(ctx context.Context, req *requests.UpdateBatchMerchantSocialRequest) ([]*db.UpdateMerchantSocialMediaLinkRow, error)
	TrashSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreSocialLink(ctx context.Context, socialID int) (bool, error)
	DeletePermanentSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreAllSocialLinks(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
