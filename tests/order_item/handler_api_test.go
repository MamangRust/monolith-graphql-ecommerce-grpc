package order_item_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderItemApiTestSuite struct {
	tests.BaseTestSuite
	gql     http.Handler
	orderID int
}

func (s *OrderItemApiTestSuite) SetupSuite() {
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
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	s.orderID = s.SeedOrder(ctx, userID, merchID, prodID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *OrderItemApiTestSuite) TestOrderItemApiQueries() {
	// The GraphQL gateway exposes order-item reads only; order items are
	// created as part of an order (see SeedOrder above).
	// 1. FindAll
	body, err := tests.DoGraphQL(s.gql, `query { findAllOrderItems(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllOrderItems"))

	// 2. FindByOrder
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findOrderItemsByOrder(input: { id: %d }) { status message data { id } } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	opData := tests.GraphQLOpData(body, "findOrderItemsByOrder")
	s.Require().NotNil(opData)
	items, _ := opData["data"].([]any)
	s.Require().NotEmpty(items, "expected order items for seeded order, got: %v", body)

	// 3. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveOrderItems(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveOrderItems"))

	// 4. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedOrderItems(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedOrderItems"))
}

func TestOrderItemApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemApiTestSuite))
}
