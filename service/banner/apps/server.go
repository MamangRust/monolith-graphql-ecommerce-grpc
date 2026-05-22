package apps

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-graphql-ecommerce-banner/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-banner/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-banner/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-banner/service"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)

	observability, _ := observability.NewObservability("banner-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: observability,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterBannerQueryServiceServer(gs, h.BannerQuery)
		pb.RegisterBannerCommandServiceServer(gs, h.BannerCommand)
	}

	return srv, nil
}
