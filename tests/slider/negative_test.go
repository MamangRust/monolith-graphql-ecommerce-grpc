package slider_test

import (
	"context"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	pbslider "github.com/MamangRust/monolith-graphql-ecommerce-pb/slider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: non-existent slider must map to codes.NotFound (404), not Internal.
func (s *SliderGapiTestSuite) TestSliderGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindById(ctx, &pbslider.FindByIdSliderRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent slider must be NotFound, got %v: %s", st.Code(), st.Message())
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *SliderRepositoryTestSuite) TestSliderFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.SliderQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
