package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

func NewRepositories(DB *db.Queries, categoryQueryClient pb.CategoryQueryServiceClient, merchantQueryClient pb.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		ProductQuery:   NewProductQueryRepository(DB),
		ProductCommand: NewProductCommandRepository(DB),
		CategoryQuery:  NewCategoryQueryRepository(categoryQueryClient),
		MerchantQuery:  NewMerchantQueryRepository(merchantQueryClient),
	}
}
