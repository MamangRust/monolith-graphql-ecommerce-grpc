package merchantapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type MerchantBaseResponseMapper interface {
	ToResponseMerchant(merchant *pb.MerchantResponse) *response.MerchantResponse
	ToResponsesMerchant(merchants []*pb.MerchantResponse) []*response.MerchantResponse
}

type MerchantQueryResponseMapper interface {
	MerchantBaseResponseMapper
	ToApiResponseMerchant(pbResponse *pb.ApiResponseMerchant) *response.ApiResponseMerchant
	ToApiResponsesMerchant(pbResponse *pb.ApiResponsesMerchant) *response.ApiResponsesMerchant
	ToApiResponsePaginationMerchant(pbResponse *pb.ApiResponsePaginationMerchant) *response.ApiResponsePaginationMerchant
	ToApiResponsePaginationMerchantDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
}

type MerchantCommandResponseMapper interface {
	MerchantBaseResponseMapper
	ToApiResponseMerchant(pbResponse *pb.ApiResponseMerchant) *response.ApiResponseMerchant
	ToResponseMerchantDeleteAt(merchant *pb.MerchantResponseDeleteAt) *response.MerchantResponseDeleteAt
	ToResponsesMerchantDeleteAt(merchants []*pb.MerchantResponseDeleteAt) []*response.MerchantResponseDeleteAt
	ToApiResponseMerchantDeleteAt(pbResponse *pb.ApiResponseMerchantDeleteAt) *response.ApiResponseMerchantDeleteAt
	ToApiResponseMerchantDelete(pbResponse *pb.ApiResponseMerchantDelete) *response.ApiResponseMerchantDelete
	ToApiResponseMerchantAll(pbResponse *pb.ApiResponseMerchantAll) *response.ApiResponseMerchantAll
	ToApiResponsePaginationMerchantDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
}
