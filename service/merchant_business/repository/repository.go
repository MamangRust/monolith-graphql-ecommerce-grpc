package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantBusinessQuery   MerchantBusinessQueryRepository
	MerchantBusinessCommand MerchantBusinessCommandRepository
}

func NewRepositories(DB *db.Queries, merchantQuery pb.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(merchantQuery),
		MerchantBusinessQuery:   NewMerchantBusinessQueryRepository(DB),
		MerchantBusinessCommand: NewMerchantBusinessCommandRepository(DB),
	}
}
