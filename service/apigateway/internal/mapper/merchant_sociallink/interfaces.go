package merchant_sociallinkgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_social_link"
)

type MerchantSocialLinkGraphqlMapper interface {
	ToGraphqlResponseMerchantSocialLink(res *pb.ApiResponseMerchantSocial) *model.APIResponseMerchantSocialMediaLink
}
