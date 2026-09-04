package review_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ReviewApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	reviewID   int
	merchantID int
	prodID     int
	userID     int
}

func (s *ReviewApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupCategoryService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()
	s.SetupReviewService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.prodID = s.SeedProduct(ctx, s.merchantID, catID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *ReviewApiTestSuite) TestReviewApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createReview(input: {
			user_id: %d
			product_id: %d
			name: "Test Review"
			comment: "Great product"
			rating: 5
		}) {
			status message
			data { id comment rating }
		}
	}`, s.userID, s.prodID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createReview")
	s.Require().NotNil(data)
	s.Equal("Great product", data["comment"])
	s.reviewID = int(data["id"].(float64))
	s.Require().NotZero(s.reviewID)

	// 2. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllReviews(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllReviews"))

	// 3. FindByProduct
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findReviewsByProduct(input: { product_id: %d }) { status message data { id } } }`, s.prodID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findReviewsByProduct"))

	// 4. FindByMerchant
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findReviewsByMerchant(input: { merchant_id: %d }) { status message data { id } } }`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findReviewsByMerchant"))

	// 5. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveReviews(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveReviews"))

	// 6. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateReview(input: { review_id: %d, comment: "Updated Comment", rating: 4 }) {
			status message
			data { id comment rating }
		}
	}`, s.reviewID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateReview")
	s.Require().NotNil(data)
	s.Equal("Updated Comment", data["comment"])

	// 7. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedReview(input: { id: %d }) { status message data { id } } }`, s.reviewID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 8. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedReviews(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedReviews"))

	// 9. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreReview(input: { id: %d }) { status message data { id } } }`, s.reviewID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedReview(input: { id: %d }) { status message } }`, s.reviewID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteReviewPermanent(input: { id: %d }) { status message } }`, s.reviewID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllReviews { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 12. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllReviewsPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestReviewApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewApiTestSuite))
}
