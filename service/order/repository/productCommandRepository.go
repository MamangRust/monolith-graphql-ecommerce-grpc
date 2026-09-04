package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	product_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/product_errors"

	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
)

type productCommandRepository struct {
	client pbproduct.ProductCommandServiceClient
}

func NewProductCommandRepository(client pbproduct.ProductCommandServiceClient) *productCommandRepository {
	return &productCommandRepository{
		client: client,
	}
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error) {
	res, err := r.client.UpdateProductCountStock(ctx, &pbproduct.UpdateProductCountStockRequest{
		ProductId: int32(product_id),
		Stock:     int32(stock),
	})
	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &db.UpdateProductCountStockRow{
		ProductID:    res.Data.Id,
		CountInStock: res.Data.CountInStock,
	}, nil
}
