package product_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ProductApiTestSuite struct {
	tests.BaseTestSuite
	gql         http.Handler
	productID   int
	categoryID  int
	categoryName string
	merchantID  int
}

func (s *ProductApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	s.categoryID = s.SeedCategory(ctx)
	s.categoryName = "Seed Category"
	s.merchantID = s.SeedMerchant(ctx, userID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *ProductApiTestSuite) TestProductApiLifecycle() {
	// 1. Create (imageProduct is optional; create with JSON)
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createProduct(input: {
			merchantId: %d
			categoryId: %d
			name: "Test Product"
			description: "Test Description"
			price: 10000
			countInStock: 10
			brand: "Test Brand"
			weight: 1000
		}) {
			status message
			data { id name }
		}
	}`, s.merchantID, s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createProduct")
	s.Require().NotNil(data)
	s.Equal("Test Product", data["name"])
	s.productID = int(data["id"].(float64))
	s.Require().NotZero(s.productID)

	// 2. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findProductById(input: { id: %d }) { status message data { id name } } }`, s.productID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findProductById")
	s.Require().NotNil(data)
	s.Equal(float64(s.productID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllProducts(input: { page: 1, pageSize: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllProducts"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveProducts(input: { page: 1, pageSize: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveProducts"))

	// 5. FindByMerchant
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findProductsByMerchant(input: { merchantId: %d }) { status message data { id } } }`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findProductsByMerchant"))

	// 6. FindByCategory
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findProductsByCategory(input: { categoryName: "%s" }) { status message data { id } } }`, s.categoryName))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findProductsByCategory"))

	// 7. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateProduct(input: {
			productId: %d
			merchantId: %d
			categoryId: %d
			name: "Updated Product"
			description: "Updated Description"
			price: 15000
			countInStock: 20
			brand: "Updated Brand"
			weight: 1200
		}) {
			status message
			data { id name }
		}
	}`, s.productID, s.merchantID, s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateProduct")
	s.Require().NotNil(data)
	s.Equal("Updated Product", data["name"])

	// 8. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedProduct(input: { id: %d }) { status message data { id } } }`, s.productID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedProducts(input: { page: 1, pageSize: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedProducts"))

	// 10. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreProduct(input: { id: %d }) { status message data { id } } }`, s.productID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedProduct(input: { id: %d }) { status message } }`, s.productID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteProductPermanent(input: { id: %d }) { status message } }`, s.productID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 12. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllProducts { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 13. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllProductsPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestProductApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductApiTestSuite))
}
