package repository

import (
	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter"
	productadapter "github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter/user"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency. Callers build them with adapter.WithDependencyGuard (typically
// around resilience.NewDependencyGuard) so the repository owns the wiring.
type GuardOptions struct {
	User    []adapter.GuardOption
	Product []adapter.GuardOption
}

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

func NewRepositories(DB *db.Queries,
	userQueryClient pbuser.UserQueryServiceClient,
	productQueryClient pbproduct.ProductQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		CartQuery:    NewCartQueryRepository(DB),
		CartCommand:  NewCartCommandRepository(DB),
		UserQuery:    useradapter.New(userQueryClient, nil, g.User...),
		ProductQuery: productadapter.New(productQueryClient, nil, g.Product...),
	}
}
