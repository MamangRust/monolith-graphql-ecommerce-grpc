package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	"github.com/MamangRust/monolith-graphql-ecommerce-shipping-address/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ShippingQuery   pb.ShippingQueryServiceServer
	ShippingCommand pb.ShippingCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ShippingQuery:   NewShippingQueryHandler(deps.Service.ShippingAddressQuery, deps.Logger),
		ShippingCommand: NewShippingCommandHandler(deps.Service.ShippingAddressCommand, deps.Logger),
	}
}
