package shipping_address_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ShippingAddressApiTestSuite struct {
	tests.BaseTestSuite
	gql               http.Handler
	orderID           int
	shippingAddressID int
}

func (s *ShippingAddressApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	s.orderID = s.SeedOrder(ctx, userID, merchID, prodID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *ShippingAddressApiTestSuite) TestShippingAddressApiLifecycle() {
	// 1. FindAll
	body, err := tests.DoGraphQL(s.gql, `query { findAllShipping(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllShipping"))

	// 2. FindByOrder
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findShippingByOrder(input: { id: %d }) { status message data { id } } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "findShippingByOrder")
	s.Require().NotNil(data)
	s.shippingAddressID = int(data["id"].(float64))
	s.Require().NotZero(s.shippingAddressID)

	// 3. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findShippingById(input: { id: %d }) { status message data { id } } }`, s.shippingAddressID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findShippingById")
	s.Require().NotNil(data)
	s.Equal(float64(s.shippingAddressID), data["id"])

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveShipping(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveShipping"))

	// 5. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedShipping(input: { id: %d }) { status message data { id } } }`, s.shippingAddressID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 6. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedShipping(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedShipping"))

	// 7. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreShipping(input: { id: %d }) { status message data { id } } }`, s.shippingAddressID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 8. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedShipping(input: { id: %d }) { status message } }`, s.shippingAddressID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteShippingPermanent(input: { id: %d }) { status message } }`, s.shippingAddressID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllShipping { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllShippingPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestShippingAddressApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingAddressApiTestSuite))
}
