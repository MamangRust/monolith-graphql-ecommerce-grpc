package transaction_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-graphql-ecommerce-transaction/repository"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"

	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-graphql-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type TransactionRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *TransactionRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupShippingAddressService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(&repository.Deps{
		Db:              queries,
		UserQueryClient:       pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		MerchantQueryClient:   pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		OrderQueryClient:      pborder.NewOrderQueryServiceClient(s.Conns["order"]),
		OrderItemQueryClient:  pborder_item.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		ShippingQueryClient:   pbshipping_address.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
	})
}

func (s *TransactionRepositoryTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)

	// 2. Create Transaction
	status := "paid"
	req := &requests.CreateTransactionRequest{
		UserID:        userID,
		OrderID:       orderID,
		MerchantID:    merchID,
		PaymentMethod: "Credit Card",
		PaymentStatus: &status,
		Amount:        100000,
	}

	created, err := s.repo.TransactionCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(orderID), created.OrderID)

	// 3. Find by ID
	found, err := s.repo.TransactionQuery.FindByID(ctx, int(created.TransactionID))
	s.NoError(err)
	s.Equal(created.Amount, found.Amount)
}

func TestTransactionRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
