package merchant_detailgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_detail"
)

type MerchantDetailGraphqlMapper interface {
	ToGraphqlResponseMerchantDetailRelation(res *pb.ApiResponseMerchantDetail) *model.APIResponseMerchantDetailRelation
	ToGraphqlResponseMerchantDetailDelete(res *pbmerchant.ApiResponseMerchantDelete) *model.APIResponseMerchantDetailDelete
	ToGraphqlResponseMerchantDetailAll(res *pbmerchant.ApiResponseMerchantAll) *model.APIResponseMerchantDetailAll
	ToGraphqlResponseMerchantDetail(res *pb.ApiResponseMerchantDetail) *model.APIResponseMerchantDetail
	ToGraphqlResponseMerchantDetailDeleteAt(res *pb.ApiResponseMerchantDetailDeleteAt) *model.APIResponseMerchantDetailDeleteAt
	ToGraphqlResponsePaginationMerchantDetail(res *pb.ApiResponsePaginationMerchantDetail) *model.APIResponsePaginationMerchantDetail
	ToGraphqlResponsePaginationMerchantDetailDeleteAt(res *pb.ApiResponsePaginationMerchantDetailDeleteAt) *model.APIResponsePaginationMerchantDetailDeleteAt
}
