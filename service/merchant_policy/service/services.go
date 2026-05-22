package service

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	mencache "github.com/MamangRust/monolith-graphql-ecommerce-merchant_policy/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_policy/repository"
)

type Service struct {
	MerchantPoliciesQuery   MerchantPoliciesQueryService
	MerchantPoliciesCommand MerchantPoliciesCommandService
}

type Deps struct {
	Cache         *mencache.Mencache
	Repository    *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		MerchantPoliciesQuery: NewMerchantPoliciesQueryService(&MerchantPoliciesQueryServiceDeps{
			Observability:            deps.Observability,
			Cache:                    deps.Cache.MerchantPoliciesQueryCache,
			MerchantPolicyRepository: deps.Repository.MerchantPoliciesQuery,
			Logger:                   deps.Logger,
		}),
		MerchantPoliciesCommand: NewMerchantPoliciesCommandService(&MerchantPoliciesCommandServiceDeps{
			Observability:            deps.Observability,
			Cache:                    deps.Cache.MerchantPoliciesCommandCache,
			MerchantPolicyRepository: deps.Repository.MerchantPoliciesCommand,
			Logger:                   deps.Logger,
		}),
	}
}
