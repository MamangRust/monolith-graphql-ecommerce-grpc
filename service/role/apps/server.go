package apps

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-role/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-role/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-role/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-role/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"google.golang.org/grpc"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	obs, _ := observability.NewObservability("role-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbrole.RegisterRoleQueryServiceServer(gs, h.RoleQuery)
		pbrole.RegisterRoleCommandServiceServer(gs, h.RoleCommand)
	}

	return srv, nil
}
