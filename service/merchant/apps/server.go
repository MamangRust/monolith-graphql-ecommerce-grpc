package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-graphql-ecommerce-merchant/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_document"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := viper.GetString("GRPC_USER_ADDR")

	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)

	repos := repository.NewRepositories(srv.DB, userQueryClient)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})
	mencache := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability(viper.GetString("merchant-server"), srv.Logger)

	svc := service.NewService(&service.Deps{
		Kafka:         myKafka,
		Repositories:  repos,
		Mencache:      mencache,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbmerchant.RegisterMerchantQueryServiceServer(gs, h.MerchantQuery)
		pbmerchant.RegisterMerchantCommandServiceServer(gs, h.MerchantCommandHandler)
		pbmerchant_document.RegisterMerchantDocumentQueryServiceServer(gs, h.MerchantDocumentQuery)
		pbmerchant_document.RegisterMerchantDocumentCommandServiceServer(gs, h.MerchantDocumentCommand)
	}

	return srv, nil
}
