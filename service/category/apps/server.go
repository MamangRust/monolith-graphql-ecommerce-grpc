package apps

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-category/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-category/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-category/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-category/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"google.golang.org/grpc"

	pbcategory "github.com/MamangRust/monolith-graphql-ecommerce-pb/category"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	obs, _ := observability.NewObservability("category-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbcategory.RegisterCategoryQueryServiceServer(gs, h.CategoryQuery)
		pbcategory.RegisterCategoryCommandServiceServer(gs, h.CategoryCommand)
		pbcategory.RegisterCategoryStatsServiceServer(gs, h.CategoryStats)
		pbcategory.RegisterCategoryStatsByIdServiceServer(gs, h.CategoryStatsById)
		pbcategory.RegisterCategoryStatsByMerchantServiceServer(gs, h.CategoryStatsByMerchant)
	}

	return srv, nil
}
