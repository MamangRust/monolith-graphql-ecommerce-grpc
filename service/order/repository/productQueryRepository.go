package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	product_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/product_errors"

	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
)

type productQueryRepository struct {
	client pbproduct.ProductQueryServiceClient
}

func NewProductQueryRepository(client pbproduct.ProductQueryServiceClient) *productQueryRepository {
	return &productQueryRepository{
		client: client,
	}
}

func (r *productQueryRepository) FindByID(ctx context.Context, product_id int) (*db.GetProductByIDRow, error) {
	res, err := r.client.FindById(ctx, &pbproduct.FindByIdProductRequest{Id: int32(product_id)})
	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &db.GetProductByIDRow{
		ProductID:   res.Data.Id,
		Name:        res.Data.Name,
		Description: &res.Data.Description,
		Price:       res.Data.Price,
		CountInStock: res.Data.CountInStock,
	}, nil
}
