package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	merchant_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/merchant"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
)

type merchantQueryRepository struct {
	client pbmerchant.MerchantQueryServiceClient
}

func NewMerchantQueryRepository(client pbmerchant.MerchantQueryServiceClient) *merchantQueryRepository {
	return &merchantQueryRepository{
		client: client,
	}
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, merchant_id int) (*db.GetMerchantByIDRow, error) {
	res, err := r.client.FindById(ctx, &pbmerchant.FindByIdMerchantRequest{Id: int32(merchant_id)})
	if err != nil {
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}

	return &db.GetMerchantByIDRow{
		MerchantID:   res.Data.Id,
		UserID:       res.Data.UserId,
		Name:         res.Data.Name,
		Description:  &res.Data.Description,
		Address:      &res.Data.Address,
		ContactEmail: &res.Data.ContactEmail,
		ContactPhone: &res.Data.ContactPhone,
		Status:       res.Data.Status,
	}, nil
}
