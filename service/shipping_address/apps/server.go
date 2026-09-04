package apps

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-shipping-address/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-shipping-address/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-shipping-address/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-shipping-address/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"google.golang.org/grpc"

	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	mencache := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability("shipping_address-server", srv.Logger)

	svc := service.NewService(&service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbshipping_address.RegisterShippingQueryServiceServer(gs, h.ShippingQuery)
		pbshipping_address.RegisterShippingCommandServiceServer(gs, h.ShippingCommand)
	}

	return srv, nil
}
