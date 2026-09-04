package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	order_item_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/order_item_errors"

	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
)

type orderItemQueryRepository struct {
	queryClient   pborder_item.OrderItemQueryServiceClient
	commandClient pborder_item.OrderItemCommandServiceClient
}

func NewOrderItemQueryRepository(queryClient pborder_item.OrderItemQueryServiceClient, commandClient pborder_item.OrderItemCommandServiceClient) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		queryClient:   queryClient,
		commandClient: commandClient,
	}
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.queryClient.FindOrderItemByOrder(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return nil, order_item_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}

	var items []*db.GetOrderItemsByOrderRow
	for _, item := range res.Data {
		items = append(items, &db.GetOrderItemsByOrderRow{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	return items, nil
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, order_id int) (*int32, error) {
	res, err := r.commandClient.CalculateTotalPrice(ctx, &pborder_item.CalculateTotalPriceRequest{OrderId: int32(order_id)})
	if err != nil {
		return nil, order_item_errors.ErrCalculateTotalPrice.WithInternal(err)
	}

	total := int32(res.TotalPrice)
	return &total, nil
}
