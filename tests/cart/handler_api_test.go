package cart_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CartApiTestSuite struct {
	tests.BaseTestSuite
	gql    http.Handler
	userID int
}

func (s *CartApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupCartService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *CartApiTestSuite) addToCart(productID int, quantity int) map[string]any {
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createCart(input: { user_id: %d, product_id: %d, quantity: %d }) {
			status message
			data { id product_id quantity }
		}
	}`, s.userID, productID, quantity))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	return tests.GraphQLOpEnvelopeData(body, "createCart")
}

func (s *CartApiTestSuite) findAllCartIDs() []any {
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findAllCarts(input: { user_id: %d, page: 1, page_size: 10 }) {
			status message
			data { id product_id quantity }
		}
	}`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	opData := tests.GraphQLOpData(body, "findAllCarts")
	s.Require().NotNil(opData)
	items, _ := opData["data"].([]any)
	return items
}

func (s *CartApiTestSuite) TestCartApiLifecycle() {
	ctx := context.Background()
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, s.userID)
	prod1ID := s.SeedProduct(ctx, merchantID, categoryID)
	prod2ID := s.SeedProduct(ctx, merchantID, categoryID)

	// 1. Create two cart items
	cart1 := s.addToCart(prod1ID, 2)
	s.Require().NotNil(cart1)
	s.Equal(float64(prod1ID), cart1["product_id"])
	s.Equal(float64(2), cart1["quantity"])
	cart1ID := int(cart1["id"].(float64))

	cart2 := s.addToCart(prod2ID, 1)
	s.Require().NotNil(cart2)
	cart2ID := int(cart2["id"].(float64))

	// 2. FindAll (both items present)
	items := s.findAllCartIDs()
	s.Require().Len(items, 2, "expected two cart items, got %d", len(items))

	// 3. Delete the first item
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		deleteCart(input: { cart_id: %d, user_id: %d }) { status message }
	}`, cart1ID, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 4. DeleteAll (the remaining item)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		deleteAllCarts(input: { user_id: %d, cart_ids: [%d] }) { status message }
	}`, s.userID, cart2ID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 5. FindAll again (must be empty)
	items = s.findAllCartIDs()
	s.Require().Empty(items, "expected no cart items after deletes, got %d", len(items))
}

func TestCartApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartApiTestSuite))
}
