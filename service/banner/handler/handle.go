package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-banner/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	BannerQuery   BannerQueryHandler
	BannerCommand BannerCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		BannerQuery:   NewBannerQueryHandler(deps.Service.BannerQuery, deps.Logger),
		BannerCommand: NewBannerCommandHandler(deps.Service.BannerCommand, deps.Logger),
	}
}
