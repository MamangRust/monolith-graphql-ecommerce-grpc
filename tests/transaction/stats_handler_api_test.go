package transaction_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	merchantID int
}

func (s *TransactionStatsApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupShippingAddressService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, userID, s.merchantID, prodID)

	// Seed a successful transaction
	_, err := s.DBPool().Exec(ctx, `
		INSERT INTO transactions (merchant_id, order_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		s.merchantID, orderID, 100000, "credit_card", "success", time.Now())
	s.Require().NoError(err)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccess() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthStatusSuccess(input: { year: %d, month: %d }) { status message }
	}`, now.Year(), int(now.Month())))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthStatusSuccess"))
}

func (s *TransactionStatsApiTestSuite) TestFindYearStatusSuccess() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findYearStatusSuccess(input: { year: %d, month: %d }) { status message }
	}`, now.Year(), int(now.Month())))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findYearStatusSuccess"))
}

func (s *TransactionStatsApiTestSuite) TestFindMonthMethodSuccess() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthMethodSuccess(input: { year: %d, month: %d }) { status message }
	}`, now.Year(), int(now.Month())))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthMethodSuccess"))
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccessByMerchant() {
	now := time.Now()
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findMonthStatusSuccessByMerchant(input: { year: %d, month: %d, merchant_id: %d }) { status message }
	}`, now.Year(), int(now.Month()), s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthStatusSuccessByMerchant"))
}

func TestTransactionStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsApiTestSuite))
}
