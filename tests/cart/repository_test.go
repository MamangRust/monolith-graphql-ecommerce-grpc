package cart_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-graphql-ecommerce-cart/repository"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"

	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type CartRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *CartRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(
		queries,
		pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		pbproduct.NewProductQueryServiceClient(s.Conns["product"]),
	)
}

func (s *CartRepositoryTestSuite) TestCartLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	childUserID := s.SeedUser(ctx)
	childCategoryID := s.SeedCategory(ctx)
	childMerchantID := s.SeedMerchant(ctx, childUserID)
	prodID := s.SeedProduct(ctx, childMerchantID, childCategoryID)

	// Fetch product details for the cart record
	prodRes, err := pbproduct.NewProductQueryServiceClient(s.Conns["product"]).FindById(ctx, &pbproduct.FindByIdProductRequest{Id: int32(prodID)})
	s.Require().NoError(err)
	s.Require().NotNil(prodRes)
	s.Require().NotNil(prodRes.Data)

	// 2. Add to Cart
	req := &requests.CartCreateRecord{
		UserID:       childUserID,
		ProductID:    prodID,
		Name:         prodRes.Data.Name,
		Price:        int(prodRes.Data.Price),
		ImageProduct: prodRes.Data.ImageProduct,
		Quantity:     2,
		Weight:       int(prodRes.Data.Weight),
	}

	created, err := s.repo.CartCommand.CreateCart(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(childUserID), created.UserID)

	// 3. Find by User ID
	items, err := s.repo.CartQuery.FindCarts(ctx, &requests.FindAllCarts{
		UserID:   childUserID,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(items)
}

func TestCartRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartRepositoryTestSuite))
}
