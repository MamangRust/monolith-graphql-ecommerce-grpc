package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	"github.com/MamangRust/monolith-graphql-ecommerce-review-detail/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewDetail        ReviewDetailHandleGrpc
	ReviewDetailQuery   pb.ReviewDetailQueryServiceServer
	ReviewDetailCommand pb.ReviewDetailCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ReviewDetailQuery:   NewReviewDetailQueryHandler(deps.Service.ReviewDetailQuery, deps.Logger),
		ReviewDetailCommand: NewReviewDetailCommandHandler(deps.Service.ReviewDetailCommand, deps.Logger),
	}
}
