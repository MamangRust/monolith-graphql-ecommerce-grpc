package merchant_social_link_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcMerchantSocialLinkNotFound  = errors.NewGrpcError("MerchantSocialLink not found", int(codes.NotFound))
	ErrGrpcMerchantSocialLinkInvalidId = errors.NewGrpcError("Invalid MerchantSocialLink ID", int(codes.NotFound))

	ErrGrpcValidateCreateMerchantSocialLink = errors.NewGrpcError("validation failed: invalid create MerchantSocialLink request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantSocialLink = errors.NewGrpcError("validation failed: invalid update MerchantSocialLink request", int(codes.InvalidArgument))
)
