package order_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-graphql-ecommerce-order/repository"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	pbtransaction "github.com/MamangRust/monolith-graphql-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type OrderRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(&repository.Deps{
		DB:               queries,
		MerchantQuery:    pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		ProductQuery:     pbproduct.NewProductQueryServiceClient(s.Conns["product"]),
		ProductCommand:   pbproduct.NewProductCommandServiceClient(s.Conns["product"]),
		OrderItemQuery:   pborder_item.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		OrderItemCommand: pborder_item.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
		UserQuery:        pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		ShippingCommand:  pbshipping_address.NewShippingCommandServiceClient(s.Conns["shipping-address"]),
		TransactionCommand: pbtransaction.NewTransactionCommandServiceClient(s.Conns["transaction"]),
	})
}

func (s *OrderRepositoryTestSuite) TestOrderLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	s.SeedProduct(ctx, merchID, catID)

	req := &requests.CreateOrderRecordRequest{
		UserID:     userID,
		MerchantID: merchID,
		TotalPrice: 5000,
	}

	created, err := s.repo.OrderCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(1), created.UserID)

	// 2. Find by ID
	found, err := s.repo.OrderQuery.FindByID(ctx, int(created.OrderID))
	s.NoError(err)
	s.Equal(int32(5000), found.TotalPrice)
}

func TestOrderRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderRepositoryTestSuite))
}
