package merchant_businessgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_business"
)

type MerchantBusinessGraphqlMapper interface {
	ToGraphqlResponseMerchantBusinessDelete(res *pbmerchant.ApiResponseMerchantDelete) *model.APIResponseMerchantBusinessDelete
	ToGraphqlResponseMerchantBusinessAll(res *pbmerchant.ApiResponseMerchantAll) *model.APIResponseMerchantBusinessAll
	ToGraphqlResponseMerchantBusiness(res *pb.ApiResponseMerchantBusiness) *model.APIResponseMerchantBusiness
	ToGraphqlResponseMerchantBusinessDeleteAt(res *pb.ApiResponseMerchantBusinessDeleteAt) *model.APIResponseMerchantBusinessDeleteAt
	ToGraphqlResponsesMerchantBusiness(res *pb.ApiResponsesMerchantBusiness) *model.APIResponsesMerchantBusiness
	ToGraphqlResponsePaginationMerchantBusinessDeleteAt(res *pb.ApiResponsePaginationMerchantBusinessDeleteAt) *model.APIResponsePaginationMerchantBusinessDeleteAt
	ToGraphqlResponsePaginationMerchantBusiness(res *pb.ApiResponsePaginationMerchantBusiness) *model.APIResponsePaginationMerchantBusiness
}
