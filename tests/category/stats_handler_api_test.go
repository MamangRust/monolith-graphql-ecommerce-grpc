package category_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CategoryStatsApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	categoryID int
	merchantID int
	userID     int
}

func (s *CategoryStatsApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()
	s.SetupOrderService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.categoryID = s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, s.categoryID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	_, err := s.DBPool().Exec(ctx, "UPDATE orders SET created_at = $1 WHERE order_id = $2",
		time.Now(), orderID)
	s.Require().NoError(err)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *CategoryStatsApiTestSuite) TestFindMonthTotalPrice() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthlyTotalPrices(input: { year: %d, month: %d }) { status message }
	}`, now.Year(), int(now.Month())))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthlyTotalPrices"))
}

func (s *CategoryStatsApiTestSuite) TestFindYearTotalPrice() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findYearlyTotalPrices(input: { year: %d }) { status message }
	}`, now.Year()))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findYearlyTotalPrices"))
}

func (s *CategoryStatsApiTestSuite) TestFindMonthPrice() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthPrice(input: { year: %d }) { status message }
	}`, now.Year()))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthPrice"))
}

func (s *CategoryStatsApiTestSuite) TestFindYearPrice() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findYearPrice(input: { year: %d }) { status message }
	}`, now.Year()))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findYearPrice"))
}

func (s *CategoryStatsApiTestSuite) TestFindMonthTotalPriceById() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthlyTotalPricesById(input: { year: %d, month: %d, category_id: %d }) { status message }
	}`, now.Year(), int(now.Month()), s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthlyTotalPricesById"))
}

func (s *CategoryStatsApiTestSuite) TestFindMonthTotalPriceByMerchant() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthlyTotalPricesByMerchant(input: { year: %d, month: %d, merchant_id: %d }) { status message }
	}`, now.Year(), int(now.Month()), s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthlyTotalPricesByMerchant"))
}

func TestCategoryStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsApiTestSuite))
}
