package merchant_awardgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_award"
)

type MerchantAwardGraphqlMapper interface {
	ToGraphqlResponseMerchantAwardDelete(res *pbmerchant.ApiResponseMerchantDelete) *model.APIResponseMerchantAwardDelete
	ToGraphqlResponseMerchantAwardAll(res *pbmerchant.ApiResponseMerchantAll) *model.APIResponseMerchantAwardAll
	ToGraphqlResponseMerchantAward(res *pb.ApiResponseMerchantAward) *model.APIResponseMerchantAward
	ToGraphqlResponseMerchantAwardDeleteAt(res *pb.ApiResponseMerchantAwardDeleteAt) *model.APIResponseMerchantAwardDeleteAt
	ToGraphqlResponseMerchantAwards(res *pb.ApiResponsesMerchantAward) *model.APIResponsesMerchantAward
	ToGraphqlResponsePaginationMerchantAwardDeleteAt(res *pb.ApiResponsePaginationMerchantAwardDeleteAt) *model.APIResponsePaginationMerchantAwardDeleteAt
	ToGraphqlPaginationMerchantAward(res *pb.ApiResponsePaginationMerchantAward) *model.APIResponsePaginationMerchantAward
}
