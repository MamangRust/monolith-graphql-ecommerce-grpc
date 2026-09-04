package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbcategory "github.com/MamangRust/monolith-graphql-ecommerce-pb/category"
	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
)

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

func NewRepositories(DB *db.Queries, categoryQueryClient pbcategory.CategoryQueryServiceClient, merchantQueryClient pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		ProductQuery:   NewProductQueryRepository(DB),
		ProductCommand: NewProductCommandRepository(DB),
		CategoryQuery:  NewCategoryQueryRepository(categoryQueryClient),
		MerchantQuery:  NewMerchantQueryRepository(merchantQueryClient),
	}
}
