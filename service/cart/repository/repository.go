package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

func NewRepositories(DB *db.Queries,
	userQuery pbuser.UserQueryServiceClient,
	productQuery pbproduct.ProductQueryServiceClient,
) *Repositories {
	return &Repositories{
		CartQuery:    NewCartQueryRepository(DB),
		CartCommand:  NewCartCommandRepository(DB),
		UserQuery:    NewUserQueryRepository(userQuery),
		ProductQuery: NewProductQueryRepository(productQuery),
	}
}
