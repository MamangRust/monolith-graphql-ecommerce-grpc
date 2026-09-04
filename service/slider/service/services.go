package service

import (
	mencache "github.com/MamangRust/monolith-graphql-ecommerce-slider/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-slider/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
)

type Service struct {
	SliderQuery   SliderQueryService
	SliderCommand SliderCommandService
}

type Deps struct {
	Repositories  *repository.Repositories
	Mencache      mencache.SliderMencache
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		SliderQuery: NewSliderQueryService(&SliderQueryServiceDeps{
			Repositories:  deps.Repositories.SliderQuery,
			Cache:         deps.Mencache,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		SliderCommand: NewSliderCommandService(&SliderCommandServiceDeps{
			Repositories:  deps.Repositories.SliderCommand,
			Cache:         deps.Mencache,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
	}
}
