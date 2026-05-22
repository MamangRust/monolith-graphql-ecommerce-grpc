package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_business/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantBusinessQuery   MerchantBusinessQueryHandler
	MerchantBusinessCommand MerchantBusinessCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantBusinessQuery:   NewMerchantBusinessQueryHandler(deps.Service.MerchantBusinessQuery, deps.Logger),
		MerchantBusinessCommand: NewMerchantBusinessCommandHandler(deps.Service.MerchantBusinessCommand, deps.Logger),
	}
}
