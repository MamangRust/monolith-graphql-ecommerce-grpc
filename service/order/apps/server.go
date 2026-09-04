package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-graphql-ecommerce-order/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-order/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-order/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-order/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-graphql-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
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

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")
	
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	productQueryClient := pbproduct.NewProductQueryServiceClient(productConn)
	productCommandClient := pbproduct.NewProductCommandServiceClient(productConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")
	if merchantAddr == "" {
		merchantAddr = "localhost:50055"
	}
	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	orderItemAddr := viper.GetString("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "localhost:50056"
	}
	orderItemConn, err := grpc.NewClient(orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order_item service: %w", err)
	}
	orderItemQueryClient := pborder_item.NewOrderItemQueryServiceClient(orderItemConn)
	orderItemCommandClient := pborder_item.NewOrderItemCommandServiceClient(orderItemConn)

	shippingAddr := viper.GetString("GRPC_SHIPPING_ADDRESS_ADDR")
	if shippingAddr == "" {
		shippingAddr = "localhost:50063"
	}
	shippingConn, err := grpc.NewClient(shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping_address service: %w", err)
	}
	shippingCommandClient := pbshipping_address.NewShippingCommandServiceClient(shippingConn)
	shippingQueryClient := pbshipping_address.NewShippingQueryServiceClient(shippingConn)

	transactionAddr := viper.GetString("GRPC_TRANSACTION_ADDR")
	if transactionAddr == "" {
		transactionAddr = "localhost:50061"
	}
	transactionConn, err := grpc.NewClient(transactionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to transaction service: %w", err)
	}
	transactionCommandClient := pbtransaction.NewTransactionCommandServiceClient(transactionConn)

	repos := repository.NewRepositories(&repository.Deps{
		DB:               srv.DB,
		UserQuery:        userQueryClient,
		ProductQuery:     productQueryClient,
		ProductCommand:   productCommandClient,
		MerchantQuery:    merchantQueryClient,
		OrderItemQuery:   orderItemQueryClient,
		OrderItemCommand: orderItemCommandClient,		ShippingCommand:  shippingCommandClient,
		ShippingQuery:    shippingQueryClient,
		TransactionCommand: transactionCommandClient,
	})

	obs, _ := observability.NewObservability("order-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})
	
	srv.RegisterServices = func(gs *grpc.Server) {
		pborder.RegisterOrderQueryServiceServer(gs, h.OrderQuery)
		pborder.RegisterOrderStatsServiceServer(gs, h.OrderStats)
		pborder.RegisterOrderCommandServiceServer(gs, h.OrderCommand)
	}

	return srv, nil
}
