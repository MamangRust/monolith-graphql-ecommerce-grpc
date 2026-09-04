package apps

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-slider/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-slider/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-slider/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-slider/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"google.golang.org/grpc"

	pbslider "github.com/MamangRust/monolith-graphql-ecommerce-pb/slider"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	cache := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability("slider-server", srv.Logger)

	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Mencache:      cache,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbslider.RegisterSliderQueryServiceServer(gs, h.SliderQuery)
		pbslider.RegisterSliderCommandServiceServer(gs, h.SliderCommand)
	}

	return srv, nil
}
