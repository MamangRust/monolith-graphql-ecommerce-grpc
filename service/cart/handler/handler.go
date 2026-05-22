package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-cart/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	CartQuery   CartQueryHandler
	CartCommand CartCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		CartQuery:   NewCartQueryHandler(deps.Service.CartQuery, deps.Logger),
		CartCommand: NewCartCommandHandler(deps.Service.CartCommand, deps.Logger),
	}
}
