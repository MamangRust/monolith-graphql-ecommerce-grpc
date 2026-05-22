package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	"github.com/MamangRust/monolith-graphql-ecommerce-slider/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	SliderQuery   pb.SliderQueryServiceServer
	SliderCommand pb.SliderCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		SliderQuery:   NewSliderQueryHandler(deps.Service.SliderQuery, deps.Logger),
		SliderCommand: NewSliderCommandHandler(deps.Service.SliderCommand, deps.Logger),
	}
}
