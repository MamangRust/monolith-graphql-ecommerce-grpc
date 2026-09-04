package review_detail_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ReviewDetailApiTestSuite struct {
	tests.BaseTestSuite
	gql            http.Handler
	reviewDetailID int
	reviewID       int
}

func (s *ReviewDetailApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()
	s.SetupReviewService()
	s.SetupReviewDetailService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	s.SeedOrder(ctx, userID, merchID, prodID)
	s.reviewID = s.SeedReview(ctx, userID, prodID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *ReviewDetailApiTestSuite) TestReviewDetailApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createReviewDetail(input: { reviewId: %d, type: "photo", url: "http://example.com/img.jpg", caption: "product photo" }) {
			status message
			data { id reviewId }
		}
	}`, s.reviewID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createReviewDetail")
	s.Require().NotNil(data)
	s.reviewDetailID = int(data["id"].(float64))
	s.Require().NotZero(s.reviewDetailID)

	// 2. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findReviewDetailById(input: { id: %d }) { status message data { id } } }`, s.reviewDetailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findReviewDetailById")
	s.Require().NotNil(data)
	s.Equal(float64(s.reviewDetailID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllReviewDetails(input: { page: 1, pageSize: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllReviewDetails"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveReviewDetails(input: { page: 1, pageSize: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveReviewDetails"))

	// 5. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateReviewDetail(input: { reviewDetailId: %d, type: "video", url: "http://example.com/vid.mp4", caption: "updated caption" }) {
			status message
			data { id }
		}
	}`, s.reviewDetailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedReviewDetail(input: { id: %d }) { status message data { id } } }`, s.reviewDetailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedReviewDetails(input: { page: 1, pageSize: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedReviewDetails"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreReviewDetail(input: { id: %d }) { status message data { id } } }`, s.reviewDetailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedReviewDetail(input: { id: %d }) { status message } }`, s.reviewDetailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteReviewDetailPermanent(input: { id: %d }) { status message } }`, s.reviewDetailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllReviewDetails { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllReviewDetailsPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestReviewDetailApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailApiTestSuite))
}
