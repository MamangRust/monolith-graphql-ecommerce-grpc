package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-graphql-ecommerce-transaction/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-transaction/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-transaction/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-transaction/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-graphql-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	pbtransaction "github.com/MamangRust/monolith-graphql-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// gRPC Client Connections
	userAddr := viper.GetString("GRPC_USER_ADDR")

	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	orderAddr := viper.GetString("GRPC_ORDER_ADDR")

	orderConn, err := grpc.NewClient(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}
	orderQueryClient := pborder.NewOrderQueryServiceClient(orderConn)

	orderItemAddr := viper.GetString("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "localhost:50056"
	}
	orderItemConn, err := grpc.NewClient(orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order_item service: %w", err)
	}
	orderItemQueryClient := pborder_item.NewOrderItemQueryServiceClient(orderItemConn)

	shippingAddr := viper.GetString("GRPC_SHIPPING_ADDRESS_ADDR")
	if shippingAddr == "" {
		shippingAddr = "localhost:50063"
	}
	shippingConn, err := grpc.NewClient(shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping_address service: %w", err)
	}
	shippingQueryClient := pbshipping_address.NewShippingQueryServiceClient(shippingConn)

	repos := repository.NewRepositories(&repository.Deps{
		DB:             srv.DB,
		UserQuery:      userQueryClient,
		MerchantQuery:  merchantQueryClient,
		OrderQuery:     orderQueryClient,
		OrderItemQuery: orderItemQueryClient,
		ShippingQuery:  shippingQueryClient,
	})
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})
	obs, _ := observability.NewObservability("transaction-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Kafka:         myKafka,
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbtransaction.RegisterTransactionQueryServiceServer(gs, h.TransactionQuery)
		pbtransaction.RegisterTransactionCommandServiceServer(gs, h.TransactionCommand)
		pbtransaction.RegisterTransactionStatsServiceServer(gs, h.TransactionStats)
		pbtransaction.RegisterTransactionStatsByMerchantServiceServer(gs, h.TransactionStatsByMerchant)
	}

	return srv, nil
}
