package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	MerchantQuery             MerchantQueryRepository
	MerchantDetailQuery       MerchantDetailQueryRepository
	MerchantDetailCommand     MerchantDetailCommandRepository
	MerchantSocialLinkCommand MerchantSocialLinkCommandRepository
}

func NewRepositories(db *db.Queries, merchantQuery pb.MerchantQueryServiceClient, pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		MerchantQuery:             NewMerchantQueryRepository(merchantQuery),
		MerchantDetailQuery:       NewMerchantDetailQueryRepository(db),
		MerchantDetailCommand:     NewMerchantDetailCommandRepository(db),
		MerchantSocialLinkCommand: NewMerchantSocialLinkCommandRepository(pool, db),
	}
}
