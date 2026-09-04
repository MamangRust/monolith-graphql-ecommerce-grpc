package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	orderitem_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/order_item_errors"

	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
)

type orderItemRepository struct {
	client pborder_item.OrderItemQueryServiceClient
}

func NewOrderItemRepository(client pborder_item.OrderItemQueryServiceClient) *orderItemRepository {
	return &orderItemRepository{
		client: client,
	}
}

func (r *orderItemRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.client.FindOrderItemByOrder(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}

	var items []*db.GetOrderItemsByOrderRow
	for _, item := range res.Data {
		items = append(items, &db.GetOrderItemsByOrderRow{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    int32(item.Quantity),
			Price:       int32(item.Price),
		})
	}

	return items, nil
}
