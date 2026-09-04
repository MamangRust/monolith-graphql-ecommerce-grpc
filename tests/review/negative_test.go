package review_test

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	pbreview "github.com/MamangRust/monolith-graphql-ecommerce-pb/review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: review has no FindById query RPC; Update on a non-existent review must
// map to codes.NotFound (404), not Internal.
func (s *ReviewGapiTestSuite) TestReviewGapiNotFound() {
	ctx := context.Background()
	_, err := s.commandClient.Update(ctx, &pbreview.UpdateReviewRequest{
		ReviewId: 999999,
		Rating:   5,
		Comment:  "not found",
	})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "update on non-existent review must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: update on a non-existent review must surface as an error (NotFound from
// the service), not a silent success.
func (s *ReviewApiTestSuite) TestReviewApiNotFound() {
	body, err := tests.DoGraphQL(s.gql, `mutation {
		updateReview(input: { review_id: 999999, rating: 5, comment: "not found" }) {
			status message
			data { id }
		}
	}`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "update on non-existent review must return a GraphQL error, got: %v", body)
	s.Require().Contains(strings.ToLower(graphqlErrors[0]), "not found", "expected a not-found style error, got: %s", graphqlErrors[0])
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *ReviewRepositoryTestSuite) TestReviewFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.ReviewQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
