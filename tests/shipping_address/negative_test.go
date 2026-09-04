package shipping_address_test

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	pbshipping_address "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: non-existent shipping address must map to codes.NotFound (404), not Internal.
func (s *ShippingAddressGapiTestSuite) TestShippingAddressGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindById(ctx, &pbshipping_address.FindByIdShippingRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent shipping address must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: non-existent shipping address must surface as an error (NotFound from the service), not a silent success.
func (s *ShippingAddressApiTestSuite) TestShippingAddressApiNotFound() {
	body, err := tests.DoGraphQL(s.gql, `query { findShippingById(input: { id: 999999 }) { status message data { id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "non-existent shipping address must return a GraphQL error, got: %v", body)
	s.Require().Contains(strings.ToLower(graphqlErrors[0]), "not found", "expected a not-found style error, got: %s", graphqlErrors[0])
}

func (s *ShippingAddressApiTestSuite) TestShippingAddressApiInvalidID() {
	body, err := tests.DoGraphQL(s.gql, `query { findShippingById(input: { id: 0 }) { status message data { id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "invalid shipping address ID must return a GraphQL error, got: %v", body)
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *ShippingAddressRepositoryTestSuite) TestShippingAddressFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.ShippingAddressQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
