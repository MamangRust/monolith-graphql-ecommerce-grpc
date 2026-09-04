package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
)

type Repositories struct {
	MerchantQuery             MerchantQueryRepository
	MerchantDetailQuery       MerchantDetailQueryRepository
	MerchantDetailCommand     MerchantDetailCommandRepository
	MerchantSocialLinkCommand MerchantSocialLinkCommandRepository
}

func NewRepositories(db *db.Queries, merchantQuery pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantQuery:             NewMerchantQueryRepository(merchantQuery),
		MerchantDetailQuery:       NewMerchantDetailQueryRepository(db),
		MerchantDetailCommand:     NewMerchantDetailCommandRepository(db),
		MerchantSocialLinkCommand: NewMerchantSocialLinkCommandRepository(db),
	}
}
