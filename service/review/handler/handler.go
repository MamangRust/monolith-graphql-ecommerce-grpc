package handler

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-review/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"

	pbreview "github.com/MamangRust/monolith-graphql-ecommerce-pb/review"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewQuery   pbreview.ReviewQueryServiceServer
	ReviewCommand pbreview.ReviewCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ReviewQuery:   NewReviewQueryHandler(deps.Service.ReviewQuery, deps.Logger),
		ReviewCommand: NewReviewCommandHandler(deps.Service.ReviewCommand, deps.Logger),
	}
}

type reviewHandleGrpc struct {
	// Dummy struct for mapping receiver
}
