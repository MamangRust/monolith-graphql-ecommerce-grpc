package order_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderApiTestSuite struct {
	tests.BaseTestSuite
	gql     http.Handler
	orderID int
	userID  int
	merchID int
	prodID  int
}

func (s *OrderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupCategoryService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()
	s.SetupOrderService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	s.merchID = s.SeedMerchant(ctx, s.userID)
	s.prodID = s.SeedProduct(ctx, s.merchID, catID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *OrderApiTestSuite) TestOrderApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createOrder(input: {
			merchant_id: %d
			user_id: %d
			total_price: 1000
			items: [{ product_id: %d, quantity: 1, price: 1000 }]
			shipping: {
				alamat: "Test Alamat"
				provinsi: "Test Provinsi"
				kota: "Test Kota"
				courier: "Test Courier"
				shipping_method: "Test Method"
				shipping_cost: 100
				negara: "Test Negara"
			}
		}) {
			status message
			data { id total_price }
		}
	}`, s.merchID, s.userID, s.prodID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createOrder")
	s.Require().NotNil(data)
	s.orderID = int(data["id"].(float64))
	s.Require().NotZero(s.orderID)

	// 2. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findOrderById(input: { id: %d }) { status message data { id } } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findOrderById")
	s.Require().NotNil(data)
	s.Equal(float64(s.orderID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllOrders(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllOrders"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveOrders(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveOrders"))

	// 5. Update (resolve the created order-item and shipping IDs first)
	orderItemID := s.firstOrderItemID()
	shippingID := s.firstShippingID()

	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateOrder(input: {
			order_id: %d
			user_id: %d
			total_price: 1500
			items: [{ order_item_id: %d, product_id: %d, quantity: 1, price: 1500 }]
			shipping: {
				shipping_id: %d
				alamat: "Updated Alamat"
				provinsi: "Updated Provinsi"
				kota: "Updated Kota"
				courier: "Updated Courier"
				shipping_method: "Updated Method"
				shipping_cost: 200
				negara: "Updated Negara"
			}
		}) {
			status message
			data { id total_price }
		}
	}`, s.orderID, s.userID, orderItemID, s.prodID, shippingID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateOrder")
	s.Require().NotNil(data)
	// Order service recomputes total_price as items total + shipping cost.
	s.Equal(float64(1700), data["total_price"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashOrder(input: { id: %d }) { status message data { id } } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedOrders(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedOrders"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreOrder(input: { id: %d }) { status message data { id } } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashOrder(input: { id: %d }) { status message } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteOrderPermanent(input: { id: %d }) { status message } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllOrders { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllOrdersPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

// Test12_GetMonthlyTotalRevenue exercises the order revenue stats query.
func (s *OrderApiTestSuite) Test12_GetMonthlyTotalRevenue() {
	body, err := tests.DoGraphQL(s.gql, `query { findMonthlyTotalRevenue(input: { year: 2024, month: 1 }) { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findMonthlyTotalRevenue"))
}

// Test13_GetYearlyTotalRevenue exercises the order revenue stats query.
func (s *OrderApiTestSuite) Test13_GetYearlyTotalRevenue() {
	body, err := tests.DoGraphQL(s.gql, `query { findYearlyTotalRevenue(input: { year: 2024 }) { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findYearlyTotalRevenue"))
}

func (s *OrderApiTestSuite) firstOrderItemID() int {
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findOrderItemsByOrder(input: { id: %d }) { status message data { id } } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	opData := tests.GraphQLOpData(body, "findOrderItemsByOrder")
	s.Require().NotNil(opData)
	items, _ := opData["data"].([]any)
	s.Require().NotEmpty(items)
	first, _ := items[0].(map[string]any)
	id, ok := first["id"].(float64)
	s.Require().True(ok)
	return int(id)
}

func (s *OrderApiTestSuite) firstShippingID() int {
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findShippingByOrder(input: { id: %d }) { status message data { id } } }`, s.orderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	opData := tests.GraphQLOpEnvelopeData(body, "findShippingByOrder")
	s.Require().NotNil(opData)
	id, ok := opData["id"].(float64)
	s.Require().True(ok)
	return int(id)
}

func TestOrderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderApiTestSuite))
}
