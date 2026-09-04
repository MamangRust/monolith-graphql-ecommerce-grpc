package merchantsociallinkapimapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/response"

	pbmerchant_detail "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_detail"
	pbmerchant_social_link "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_social_link"
)

type merchantSocialLinkCommandResponseMapper struct{}

func NewMerchantSocialLinkCommandResponseMapper() MerchantSocialLinkCommandResponseMapper {
	return &merchantSocialLinkCommandResponseMapper{}
}

func (m *merchantSocialLinkCommandResponseMapper) MapMerchantSocialLink(doc *pbmerchant_detail.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse {
	if doc == nil { return nil }
	return &response.MerchantSocialLinkResponse{
		ID:               int(doc.Id),
		MerchantDetailID: int(doc.MerchantDetailId),
		Platform:         doc.Platform,
		URL:              doc.Url,
	}
}

func (m *merchantSocialLinkCommandResponseMapper) ToApiResponseMerchantSocialLink(doc *pbmerchant_social_link.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink {
	return &response.ApiResponseMerchantSocialLink{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.MapMerchantSocialLink(doc.Data),
	}
}
