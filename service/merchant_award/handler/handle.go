package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_award/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantAwardQuery   MerchantAwardQueryHandler
	MerchantAwardCommand MerchantAwardCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantAwardQuery:   NewMerchantAwardQueryHandler(deps.Service.MerchantAwardQuery, deps.Logger),
		MerchantAwardCommand: NewMerchantAwardCommandHandler(deps.Service.MerchantAwardCommand, deps.Logger),
	}
}
