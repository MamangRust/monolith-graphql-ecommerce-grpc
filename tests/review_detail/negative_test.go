package review_detail_test

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	pbreview_detail "github.com/MamangRust/monolith-graphql-ecommerce-pb/review_detail"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: non-existent review detail must map to codes.NotFound (404), not Internal.
func (s *ReviewDetailGapiTestSuite) TestReviewDetailGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindById(ctx, &pbreview_detail.FindByIdReviewDetailRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent review detail must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: non-existent review detail must surface as an error (NotFound from the service), not a silent success.
func (s *ReviewDetailApiTestSuite) TestReviewDetailApiNotFound() {
	body, err := tests.DoGraphQL(s.gql, `query { findReviewDetailById(input: { id: 999999 }) { status message data { id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "non-existent review detail must return a GraphQL error, got: %v", body)
	s.Require().Contains(strings.ToLower(graphqlErrors[0]), "not found", "expected a not-found style error, got: %s", graphqlErrors[0])
}

func (s *ReviewDetailApiTestSuite) TestReviewDetailApiInvalidID() {
	body, err := tests.DoGraphQL(s.gql, `query { findReviewDetailById(input: { id: 0 }) { status message data { id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "invalid review detail ID must return a GraphQL error, got: %v", body)
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *ReviewDetailRepositoryTestSuite) TestReviewDetailFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.ReviewDetailQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
