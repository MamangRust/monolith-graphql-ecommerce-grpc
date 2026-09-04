package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_detail/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_detail/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_detail/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_detail/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pbmerchant_detail "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_detail"
	pbmerchant_social_link "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_social_link"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(
		merchantAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}

	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	repos := repository.NewRepositories(srv.DB, merchantQueryClient)
	obs, _ := observability.NewObservability("merchant-detail-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbmerchant_detail.RegisterMerchantDetailQueryServiceServer(gs, h.MerchantDetailQuery)
		pbmerchant_detail.RegisterMerchantDetailCommandServiceServer(gs, h.MerchantDetailCommand)
		pbmerchant_social_link.RegisterMerchantSocialCommandServiceServer(gs, h.MerchantSocialLinkCommand)
	}

	return srv, nil
}
