package handler

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-slider/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"

	pbslider "github.com/MamangRust/monolith-graphql-ecommerce-pb/slider"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	SliderQuery   pbslider.SliderQueryServiceServer
	SliderCommand pbslider.SliderCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		SliderQuery:   NewSliderQueryHandler(deps.Service.SliderQuery, deps.Logger),
		SliderCommand: NewSliderCommandHandler(deps.Service.SliderCommand, deps.Logger),
	}
}
