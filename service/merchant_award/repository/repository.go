package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
)

type Repositories struct {
	MerchantAwardQuery   MerchantAwardQueryRepository
	MerchantAwardCommand MerchantAwardCommandRepository
	MerchantQuery        MerchantQueryRepository
}

func NewRepositories(db *db.Queries, merchantQuery pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantAwardQuery:   NewMerchantAwardQueryRepository(db),
		MerchantAwardCommand: NewMerchantAwardCommandRepository(db),
		MerchantQuery:        NewMerchantQueryRepository(merchantQuery),
	}
}
