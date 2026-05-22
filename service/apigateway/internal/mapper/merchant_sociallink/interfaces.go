package merchant_sociallinkgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type MerchantSocialLinkGraphqlMapper interface {
	ToGraphqlResponseMerchantSocialLink(res *pb.ApiResponseMerchantSocialMediaLink) *model.APIResponseMerchantSocialMediaLink
}
