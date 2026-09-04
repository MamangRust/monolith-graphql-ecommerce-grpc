package merchant_sociallinkgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_social_link"
	pbmd "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_detail"
)

type merchantSocialLinkResponseMapper struct {
}

func NewMerchantSocialLinkResponseMapper() *merchantSocialLinkResponseMapper {
	return &merchantSocialLinkResponseMapper{}
}

func (m *merchantSocialLinkResponseMapper) ToGraphqlResponseMerchantSocialLink(res *pb.ApiResponseMerchantSocial) *model.APIResponseMerchantSocialMediaLink {
	var data []*model.MerchantSocialMediaLinkResponse
	if res.Data != nil {
		data = []*model.MerchantSocialMediaLinkResponse{m.mapResponseMerchantSocialLink(res.Data)}
	}

	return &model.APIResponseMerchantSocialMediaLink{
		Status:  res.Status,
		Message: res.Message,
		Data:    data,
	}
}

func (m *merchantSocialLinkResponseMapper) mapResponseMerchantSocialLink(response *pbmd.MerchantSocialMediaLinkResponse) *model.MerchantSocialMediaLinkResponse {
	return &model.MerchantSocialMediaLinkResponse{
		ID:       int32(response.Id),
		Platform: response.Platform,
		URL:      response.Url,
	}
}
