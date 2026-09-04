package banner_test

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	pbbanner "github.com/MamangRust/monolith-graphql-ecommerce-pb/banner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: non-existent banner must map to codes.NotFound (404), not Internal.
func (s *BannerGapiTestSuite) TestBannerGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindById(ctx, &pbbanner.FindByIdBannerRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent banner must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: non-existent banner must surface as an error (NotFound from the service), not a silent success.
func (s *BannerApiTestSuite) TestBannerApiNotFound() {
	body, err := tests.DoGraphQL(s.gql, `query { findBannerById(input: { id: 999999 }) { status message data { banner_id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "non-existent banner must return a GraphQL error, got: %v", body)
	s.Require().Contains(strings.ToLower(graphqlErrors[0]), "not found", "expected a not-found style error, got: %s", graphqlErrors[0])
}

func (s *BannerApiTestSuite) TestBannerApiInvalidID() {
	body, err := tests.DoGraphQL(s.gql, `query { findBannerById(input: { id: 0 }) { status message data { banner_id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "invalid banner ID must return a GraphQL error, got: %v", body)
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *BannerRepositoryTestSuite) TestBannerFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.BannerQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
