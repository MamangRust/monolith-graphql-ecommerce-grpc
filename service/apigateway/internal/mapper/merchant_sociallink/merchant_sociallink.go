package merchant_sociallinkgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type merchantSocialLinkResponseMapper struct {
}

func NewMerchantSocialLinkResponseMapper() *merchantSocialLinkResponseMapper {
	return &merchantSocialLinkResponseMapper{}
}

func (m *merchantSocialLinkResponseMapper) ToGraphqlResponseMerchantSocialLink(res *pb.ApiResponseMerchantSocialMediaLink) *model.APIResponseMerchantSocialMediaLink {
	return &model.APIResponseMerchantSocialMediaLink{
		Status:  res.Status,
		Message: res.Message,
		Data:    m.mapResponsesMerchantSocialLink(res.Data),
	}
}

func (m *merchantSocialLinkResponseMapper) mapResponseMerchantSocialLink(response *pb.MerchantSocialMediaLinkResponse) *model.MerchantSocialMediaLinkResponse {
	return &model.MerchantSocialMediaLinkResponse{
		ID:       int32(response.Id),
		Platform: response.Platform,
		URL:      response.Url,
	}
}

func (m *merchantSocialLinkResponseMapper) mapResponsesMerchantSocialLink(merchants []*pb.MerchantSocialMediaLinkResponse) []*model.MerchantSocialMediaLinkResponse {
	var responses []*model.MerchantSocialMediaLinkResponse

	for _, s := range merchants {
		responses = append(responses, m.mapResponseMerchantSocialLink(s))
	}

	return responses
}
