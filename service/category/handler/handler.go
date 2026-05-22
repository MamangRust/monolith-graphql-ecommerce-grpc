package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	"github.com/MamangRust/monolith-graphql-ecommerce-category/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	CategoryQuery           pb.CategoryQueryServiceServer
	CategoryCommand         pb.CategoryCommandServiceServer
	CategoryStats           CategoryStatsHandler
	CategoryStatsById       CategoryStatsByIdHandler
	CategoryStatsByMerchant CategoryStatsByMerchantHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		CategoryQuery:           NewCategoryQueryHandler(deps.Service.CategoryQuery, deps.Logger),
		CategoryCommand:         NewCategoryCommandHandler(deps.Service.CategoryCommand, deps.Logger),
		CategoryStats:           NewCategoryStatsHandler(deps.Service.CategoryStats, deps.Logger),
		CategoryStatsById:       NewCategoryStatsByIdHandler(deps.Service.CategoryStatsById, deps.Logger),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantHandler(deps.Service.CategoryStatsByMerchant, deps.Logger),
	}
}
