package review_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-graphql-ecommerce-review/repository"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"

	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type ReviewRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *ReviewRepositoryTestSuite) SetupSuite() {
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

func (s *ReviewRepositoryTestSuite) TestReviewLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// 2. Create Review
	req := &requests.CreateReviewRequest{
		UserID:    userID,
		ProductID: prodID,
		Rating:    5,
		Comment:   "Excellent product!",
	}

	created, err := s.repo.ReviewCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(5), created.Rating)

	// 3. Find by ID
	found, err := s.repo.ReviewQuery.FindByID(ctx, int(created.ReviewID))
	s.NoError(err)
	s.Equal(created.Comment, found.Comment)

	// 4. Find by Product ID
	reviews, err := s.repo.ReviewQuery.FindByProduct(ctx, &requests.FindAllReviewByProduct{
		ProductID: prodID,
		Rating:    5,
		Page:      1,
		PageSize:  10,
	})
	s.NoError(err)
	s.NotEmpty(reviews)
}

func TestReviewRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewRepositoryTestSuite))
}
