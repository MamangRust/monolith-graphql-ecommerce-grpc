package order_item_test

import (
	"context"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	pborder_item "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
)

// order_item has no single-record lookup: FindOrderItemByOrder returns an empty
// list (success) for a non-existent order, so there is no NotFound path.
//
// gapi: a non-existent order must return an empty result, not an error.
func (s *OrderItemGapiTestSuite) TestOrderItemGapiEmptyResult() {
	ctx := context.Background()
	res, err := s.queryClient.FindOrderItemByOrder(ctx, &pborder_item.FindByIdOrderItemRequest{Id: 999999})
	s.NoError(err)
	s.NotNil(res)
	s.Empty(res.Data)
}

// api: a non-existent order returns a success response with empty data.
func (s *OrderItemApiTestSuite) TestOrderItemApiEmptyResult() {
	body, err := tests.DoGraphQL(s.gql, `query { findOrderItemsByOrder(input: { id: 999999 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	opData := tests.GraphQLOpData(body, "findOrderItemsByOrder")
	s.Require().NotNil(opData)
	items, _ := opData["data"].([]any)
	s.Empty(items)
}

// repository: FindOrderItemByOrder on a non-existent order must return an empty
// result without error.
func (s *OrderItemRepositoryTestSuite) TestOrderItemFindByOrderEmpty() {
	ctx := context.Background()
	items, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, 999999)
	s.NoError(err)
	s.Empty(items)
}
