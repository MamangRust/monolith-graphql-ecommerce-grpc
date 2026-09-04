package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantBusinessQuery   MerchantBusinessQueryRepository
	MerchantBusinessCommand MerchantBusinessCommandRepository
}

func NewRepositories(DB *db.Queries, merchantQuery pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(merchantQuery),
		MerchantBusinessQuery:   NewMerchantBusinessQueryRepository(DB),
		MerchantBusinessCommand: NewMerchantBusinessCommandRepository(DB),
	}
}
