package user_test

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: non-existent user must map to codes.NotFound (404), not Internal.
func (s *UserGapiTestSuite) TestUserGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindById(ctx, &pbuser.FindByIdUserRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent user must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: non-existent user must surface as an error (NotFound from the service), not a silent success.
func (s *UserHandlerTestSuite) TestUserApiNotFound() {
	body, err := tests.DoGraphQL(s.gql, `query { findByIdUser(input: { id: 999999 }) { status message data { id email } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "non-existent user must return a GraphQL error, got: %v", body)
	s.Require().Contains(strings.ToLower(graphqlErrors[0]), "not found", "expected a not-found style error, got: %s", graphqlErrors[0])
}

func (s *UserHandlerTestSuite) TestUserApiInvalidID() {
	body, err := tests.DoGraphQL(s.gql, `query { findByIdUser(input: { id: 0 }) { status message data { id email } } }`)
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "invalid user ID must return a GraphQL error, got: %v", body)
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *UserRepositoryTestSuite) TestUserFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.UserQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}