package order_test

import (
	"context"
	"testing"

	order_cache "github.com/MamangRust/monolith-graphql-ecommerce-order/cache"
	order_handler "github.com/MamangRust/monolith-graphql-ecommerce-order/handler"
	order_repo "github.com/MamangRust/monolith-graphql-ecommerce-order/repository"
	order_service "github.com/MamangRust/monolith-graphql-ecommerce-order/service"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-graphql-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	pbtransaction "github.com/MamangRust/monolith-graphql-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type OrderGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pborder.OrderQueryServiceClient
	commandClient pborder.OrderCommandServiceClient
}

func (s *OrderGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Order dependencies
	mencache := order_cache.NewMencache(cacheStore)
	repos := order_repo.NewRepositories(&order_repo.Deps{
		Db:               queries,
		MerchantQueryClient:    pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		ProductQueryClient:     pbproduct.NewProductQueryServiceClient(s.Conns["product"]),
		ProductCommandClient:   pbproduct.NewProductCommandServiceClient(s.Conns["product"]),
		OrderItemQueryClient:   pborder_item.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		OrderItemCommandClient: pborder_item.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
		UserQueryClient:        pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		ShippingCommandClient:  pbshipping_address.NewShippingCommandServiceClient(s.Conns["shipping-address"]),
		ShippingQueryClient:    pbshipping_address.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
		TransactionCommandClient: pbtransaction.NewTransactionCommandServiceClient(s.Conns["transaction"]),
	})
	svc := order_service.NewService(&order_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := order_handler.NewHandler(&order_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pborder.RegisterOrderQueryServiceServer(server, handler.OrderQuery)
	pborder.RegisterOrderCommandServiceServer(server, handler.OrderCommand)
	
	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pborder.NewOrderQueryServiceClient(conn)
	s.commandClient = pborder.NewOrderCommandServiceClient(conn)
}

func (s *OrderGapiTestSuite) TestOrderGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pborder.CreateOrderRequest{
		UserId:     int32(userID),
		MerchantId: int32(merchID),
		TotalPrice: 10000,
		Items: []*pborder.CreateOrderItemRequest{
			{
				ProductId: int32(prodID),
				Quantity:  1,
				Price:     10000,
			},
		},
		Shipping: &pbshipping_address.CreateShippingAddressRequest{
			Alamat:          "Test Address",
			Provinsi:        "Test Province",
			Kota:            "Test City",
			Negara:          "Test Country",
			Courier:         "Test Courier",
			ShippingMethod:  "Test Method",
			ShippingCost:    1000,
		},
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	orderID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)
	s.Equal(int32(userID), getRes.Data.UserId)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pborder.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pborder.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	// Fetch terms first
	itemClient := pborder_item.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	itemsRes, err := itemClient.FindOrderItemByOrder(ctx, &pborder_item.FindByIdOrderItemRequest{Id: orderID})
	s.Require().NoError(err)
	s.NotEmpty(itemsRes.Data)
	orderItemID := itemsRes.Data[0].Id

	_, err = s.commandClient.Update(ctx, &pborder.UpdateOrderRequest{
		OrderId:    orderID,
		UserId:     int32(userID),
		TotalPrice: 15000,
		Items: []*pborder.UpdateOrderItemRequest{
			{
				OrderItemId: orderItemID,
				ProductId:   int32(prodID),
				Quantity:    1,
				Price:       15000,
			},
		},
		Shipping: &pbshipping_address.UpdateShippingAddressRequest{
			Alamat:          "Updated Address",
			Provinsi:        "Updated Province",
			Kota:            "Updated City",
			Negara:          "Updated Country",
			Courier:         "Updated Courier",
			ShippingMethod:  "Updated Method",
			ShippingCost:    1500,
		},
	})
	s.Require().NoError(err)

	// 7. Trash
	_, err = s.commandClient.TrashedOrder(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pborder.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreOrder(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedOrder(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	_, err = s.commandClient.DeleteOrderPermanent(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllOrder(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllOrderPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestOrderGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderGapiTestSuite))
}
