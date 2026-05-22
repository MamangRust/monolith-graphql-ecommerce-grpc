package merchantsociallinkapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type MerchantSocialLinkBaseResponseMapper interface {
	MapMerchantSocialLink(doc *pb.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse
	ToApiResponseMerchantSocialLink(doc *pb.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink
}

type MerchantSocialLinkQueryResponseMapper interface {
	MerchantSocialLinkBaseResponseMapper
}

type MerchantSocialLinkCommandResponseMapper interface {
	MerchantSocialLinkBaseResponseMapper
}
