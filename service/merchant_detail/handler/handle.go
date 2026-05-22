package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_detail/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantDetailQuery       MerchantDetailQueryHandler
	MerchantDetailCommand     MerchantDetailCommandHandler
	MerchantSocialLinkCommand MerchantSocialLinkCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantDetailQuery:       NewMerchantDetailQueryHandler(deps.Service.MerchantDetailQuery, deps.Logger),
		MerchantDetailCommand:     NewMerchantDetailCommandHandler(deps.Service.MerchantDetailCommand, deps.Logger),
		MerchantSocialLinkCommand: NewMerchantSocialLinkCommandHandler(deps.Service.MerchantSocialLinkCommand, deps.Service.MerchantDetailQuery, deps.Logger),
	}
}
