package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/shipping_address_errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
)

type shippingAddressCommandRepository struct {
	client pbshipping_address.ShippingCommandServiceClient
}

func NewShippingAddressCommandRepository(client pbshipping_address.ShippingCommandServiceClient) *shippingAddressCommandRepository {
	return &shippingAddressCommandRepository{
		client: client,
	}
}

func (r *shippingAddressCommandRepository) Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*db.CreateShippingAddressRow, error) {
	res, err := r.client.CreateShipping(ctx, &pbshipping_address.CreateShippingAddressRequest{
		OrderId:        int32(*request.OrderID),
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   int32(request.ShippingCost),
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress.WithInternal(err)
	}

	return &db.CreateShippingAddressRow{
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

func (r *shippingAddressCommandRepository) Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*db.UpdateShippingAddressRow, error) {
	var shippingID int32
	if request.ShippingID != nil {
		shippingID = int32(*request.ShippingID)
	}

	res, err := r.client.UpdateShipping(ctx, &pbshipping_address.UpdateShippingAddressRequest{
		ShippingId:     shippingID,
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   int32(request.ShippingCost),
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress.WithInternal(err)
	}

	return &db.UpdateShippingAddressRow{
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

func (r *shippingAddressCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error) {
	_, err := r.client.DeleteShippingByOrderPermanent(ctx, &pbshipping_address.FindByIdShippingRequest{
		Id: int32(order_id),
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := r.client.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress.WithInternal(err)
	}

	return res.Status == "success", nil
}
