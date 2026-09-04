package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	pbtransaction "github.com/MamangRust/monolith-graphql-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type Repositories struct {
	MerchantQuery        MerchantQueryRepository
	ProductQuery         ProductQueryRepository
	ProductCommand       ProductCommandRepository
	OrderItemQuery       OrderItemQueryRepository
	OrderItemCommand     OrderItemCommandRepository
	OrderQuery           OrderQueryRepository
	OrderCommand         OrderCommandRepository
	OrderStats           OrderStatsRepository
	OrderStatsByMerchant OrderStatsByMerchantRepository
	UserQuery            UserQueryRepository
	ShippingAddress      ShippingAddressCommandRepository
	TransactionCommand   TransactionCommandRepository
	ShippingQuery        pbshipping_address.ShippingQueryServiceClient
}

type Deps struct {
	DB               *db.Queries
	MerchantQuery    pbmerchant.MerchantQueryServiceClient
	ProductQuery     pbproduct.ProductQueryServiceClient
	ProductCommand   pbproduct.ProductCommandServiceClient
	OrderItemQuery   pborder_item.OrderItemQueryServiceClient
	OrderItemCommand pborder_item.OrderItemCommandServiceClient
	UserQuery        pbuser.UserQueryServiceClient
	ShippingCommand  pbshipping_address.ShippingCommandServiceClient
	TransactionCommand pbtransaction.TransactionCommandServiceClient
	ShippingQuery    pbshipping_address.ShippingQueryServiceClient
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		MerchantQuery:    NewMerchantQueryRepository(deps.MerchantQuery),
		ProductQuery:     NewProductQueryRepository(deps.ProductQuery),
		ProductCommand:   NewProductCommandRepository(deps.ProductCommand),
		OrderItemQuery:   NewOrderItemQueryRepository(deps.OrderItemQuery, deps.OrderItemCommand),
		OrderItemCommand: NewOrderItemCommandRepository(deps.OrderItemCommand),
		OrderQuery:       NewOrderQueryRepository(deps.DB),
		OrderCommand:     NewOrderCommandRepository(deps.DB),
		UserQuery:        NewUserQueryRepository(deps.UserQuery),
		ShippingAddress:  NewShippingAddressCommandRepository(deps.ShippingCommand),
		TransactionCommand: NewTransactionCommandRepository(deps.TransactionCommand),
		ShippingQuery:    deps.ShippingQuery,
		OrderStats:       NewOrderStatsRepository(deps.DB),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(
			deps.DB,
		),
	}
}
