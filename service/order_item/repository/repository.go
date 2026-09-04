package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
)

type Repositories struct {
	OrderItemQuery   OrderItemQueryRepository
	OrderItemCommand OrderItemCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		OrderItemQuery:   NewOrderItemQueryRepository(DB),
		OrderItemCommand: NewOrderItemCommandRepository(DB),
	}
}
