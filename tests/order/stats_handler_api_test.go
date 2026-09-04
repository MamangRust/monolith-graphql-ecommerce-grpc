package order_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderStatsApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	orderID    int
	userID     int
	merchantID int
}

func (s *OrderStatsApiTestSuite) SetupSuite() {
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
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	s.orderID = s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	_, err := s.DBPool().Exec(ctx, "UPDATE orders SET created_at = $1 WHERE order_id = $2",
		time.Now(), s.orderID)
	s.Require().NoError(err)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenue() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthlyTotalRevenue(input: { year: %d, month: %d }) { status message }
	}`, now.Year(), int(now.Month())))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthlyTotalRevenue"))
}

func (s *OrderStatsApiTestSuite) TestFindYearlyTotalRevenue() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findYearlyTotalRevenue(input: { year: %d }) { status message }
	}`, now.Year()))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findYearlyTotalRevenue"))
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyOrder() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthlyRevenue(input: { year: %d, month: %d }) { status message }
	}`, now.Year(), int(now.Month())))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthlyRevenue"))
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthlyTotalRevenueByMerchant(input: { year: %d, month: %d, merchant_id: %d }) { status message }
	}`, now.Year(), int(now.Month()), s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthlyTotalRevenueByMerchant"))
}

func TestOrderStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsApiTestSuite))
}
