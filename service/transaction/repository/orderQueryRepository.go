package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	order_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/order_errors"

	pborder "github.com/MamangRust/monolith-graphql-ecommerce-pb/order"
)

type orderQueryRepository struct {
	client pborder.OrderQueryServiceClient
}

func NewOrderQueryRepository(client pborder.OrderQueryServiceClient) *orderQueryRepository {
	return &orderQueryRepository{
		client: client,
	}
}

func (r *orderQueryRepository) FindByID(ctx context.Context, order_id int) (*db.GetOrderByIDRow, error) {
	res, err := r.client.FindById(ctx, &pborder.FindByIdOrderRequest{Id: int32(order_id)})
	if err != nil {
		return nil, order_errors.ErrFindById.WithInternal(err)
	}

	return &db.GetOrderByIDRow{
		OrderID:    res.Data.Id,
		UserID:     res.Data.UserId,
		MerchantID: res.Data.MerchantId,
		TotalPrice: res.Data.TotalPrice,
	}, nil
}
