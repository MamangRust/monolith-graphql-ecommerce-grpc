package repository

import (
	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter/merchant"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	Merchant []adapter.GuardOption
}

type Repositories struct {
	MerchantBusinessQuery   MerchantBusinessQueryRepository
	MerchantBusinessCommand MerchantBusinessCommandRepository
	MerchantQuery           MerchantQueryRepository
}

func NewRepositories(db *db.Queries, merchantQueryClient pbmerchant.MerchantQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantBusinessQuery:   NewMerchantBusinessQueryRepository(db),
		MerchantBusinessCommand: NewMerchantBusinessCommandRepository(db),
		MerchantQuery:           merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
