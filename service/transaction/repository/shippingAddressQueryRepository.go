package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	shippingaddress_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/shipping_address_errors"

	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
)

type shippingAddressQueryRepository struct {
	client pbshipping_address.ShippingQueryServiceClient
}

func NewShippingAddressQueryRepository(client pbshipping_address.ShippingQueryServiceClient) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{
		client: client,
	}
}

func (r *shippingAddressQueryRepository) FindByID(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, error) {
	res, err := r.client.FindByOrder(ctx, &pbshipping_address.FindByIdShippingRequest{Id: int32(order_id)})
	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder.WithInternal(err)
	}

	return &db.GetShippingAddressByOrderIDRow{
		ShippingAddressID: res.Data.Id,
		OrderID:           res.Data.OrderId,
		Alamat:            res.Data.Alamat,
		Provinsi:          res.Data.Provinsi,
		Negara:            res.Data.Negara,
		Kota:              res.Data.Kota,
		ShippingMethod:    res.Data.ShippingMethod,
		ShippingCost:      float64(res.Data.ShippingCost),
	}, nil
}
