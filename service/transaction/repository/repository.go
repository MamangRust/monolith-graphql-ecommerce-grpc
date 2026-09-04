package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-graphql-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type Repositories struct {
	TransactionCommand TransactionCommandRepository
	TransactionQuery   TransactionQueryRepository
	OrderItem          OrderItemRepository
	OrderQuery         OrderQueryRepository
	MerchantQuery      MerchantQueryRepository
	ShippingAddress    ShippingAddressQueryRepository
	TransactionStats   TransactionStatsRepository
	StatsByMerchant    TransactionStatsByMerchantRepository
	UserQuery          UserQueryRepository
}

type Deps struct {
	DB             *db.Queries
	UserQuery      pbuser.UserQueryServiceClient
	MerchantQuery  pbmerchant.MerchantQueryServiceClient
	OrderQuery     pborder.OrderQueryServiceClient
	OrderItemQuery pborder_item.OrderItemQueryServiceClient
	ShippingQuery  pbshipping_address.ShippingQueryServiceClient
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(deps.DB),
		TransactionQuery:   NewTransactionQueryRepository(deps.DB),
		OrderItem:          NewOrderItemRepository(deps.OrderItemQuery),
		OrderQuery:         NewOrderQueryRepository(deps.OrderQuery),
		MerchantQuery:      NewMerchantQueryRepository(deps.MerchantQuery),
		ShippingAddress:    NewShippingAddressQueryRepository(deps.ShippingQuery),
		TransactionStats:   NewTransactionStatsRepository(deps.DB),
		StatsByMerchant:    NewTransactionStatsByMerchantRepository(deps.DB),
		UserQuery:          NewUserQueryRepository(deps.UserQuery),
	}
}
