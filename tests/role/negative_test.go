package role_test

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: non-existent role must map to codes.NotFound (404), not Internal.
func (s *RoleGapiTestSuite) TestRoleGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindByIdRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent role must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: non-existent role must surface as an error (NotFound from the service), not a silent success.
func (s *RoleApiTestSuite) TestRoleApiNotFound() {
	body, err := tests.DoGraphQL(s.gql, `query { findByIdRole(input: { role_id: 999999 }) { status message data { id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "non-existent role must return a GraphQL error, got: %v", body)
	s.Require().Contains(strings.ToLower(graphqlErrors[0]), "not found", "expected a not-found style error, got: %s", graphqlErrors[0])
}

func (s *RoleApiTestSuite) TestRoleApiInvalidID() {
	body, err := tests.DoGraphQL(s.gql, `query { findByIdRole(input: { role_id: 0 }) { status message data { id } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "invalid role ID must return a GraphQL error, got: %v", body)
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *RoleRepositoryTestSuite) TestRoleFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.RoleQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
