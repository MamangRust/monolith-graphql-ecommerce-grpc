package transaction_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type TransactionApiTestSuite struct {
	tests.BaseTestSuite
	gql           http.Handler
	transactionID int
	orderID       int
	merchID       int
}

func (s *TransactionApiTestSuite) SetupSuite() {
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
	s.merchID = s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchID, catID)
	s.orderID = s.SeedOrder(ctx, userID, s.merchID, prodID)
	s.SeedShippingAddress(ctx, s.orderID)
	s.SeedOrderItem(ctx, s.orderID, prodID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *TransactionApiTestSuite) TestTransactionApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createTransaction(input: {
			order_id: %d
			merchant_id: %d
			payment_method: "bank_transfer"
			amount: 100000
		}) {
			status message
			data { id order_id amount }
		}
	}`, s.orderID, s.merchID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createTransaction")
	s.Require().NotNil(data)
	s.Greater(data["amount"].(float64), float64(0))
	s.transactionID = int(data["id"].(float64))
	s.Require().NotZero(s.transactionID)

	// 2. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllTransaction(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllTransaction"))

	// 3. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findTransactionById(input: { id: %d }) { status message data { id } } }`, s.transactionID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findTransactionById")
	s.Require().NotNil(data)
	s.Equal(float64(s.transactionID), data["id"])

	// 4. FindByMerchant
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findTransactionByMerchant(input: { merchant_id: %d, page: 1, page_size: 10 }) { status message data { id } } }`, s.merchID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTransactionByMerchant"))

	// 5. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findByActive(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByActive"))

	// 6. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateTransaction(input: {
			transaction_id: %d
			order_id: %d
			merchant_id: %d
			payment_method: "bank_transfer"
			amount: 100000
		}) {
			status message
			data { id amount }
		}
	}`, s.transactionID, s.orderID, s.merchID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateTransaction")
	s.Require().NotNil(data)
	s.Greater(data["amount"].(float64), float64(0))

	// 7. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedTransaction(input: { id: %d }) { status message data { id } } }`, s.transactionID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 8. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findByTrashed(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByTrashed"))

	// 9. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreTransaction(input: { id: %d }) { status message data { id } } }`, s.transactionID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedTransaction(input: { id: %d }) { status message } }`, s.transactionID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteTransactionPermanent(input: { id: %d }) { status message } }`, s.transactionID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllTransaction { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 12. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllTransactionPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

// Test13_MonthlySuccessStats exercises the monthly success stats query.
func (s *TransactionApiTestSuite) Test13_MonthlySuccessStats() {
	body, err := tests.DoGraphQL(s.gql, `query { findMonthStatusSuccess(input: { year: 2026, month: 4 }) { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthStatusSuccess"))
}

func TestTransactionApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionApiTestSuite))
}
