package merchant_policygraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_policy"
)

type MerchantPolicyGraphqlMapper interface {
	ToGraphqlResponseMerchantPolicyDelete(res *pbmerchant.ApiResponseMerchantDelete) *model.APIResponseMerchantPolicyDelete
	ToGraphqlResponseMerchantPolicyAll(res *pbmerchant.ApiResponseMerchantAll) *model.APIResponseMerchantPolicyAll
	ToGraphqlResponseMerchantPolicy(res *pb.ApiResponseMerchantPolicies) *model.APIResponseMerchantPolicy
	ToGraphqlResponseMerchantPolicyDeleteAt(res *pb.ApiResponseMerchantPoliciesDeleteAt) *model.APIResponseMerchantPolicyDeleteAt
	ToGraphqlResponsesMerchantPolicy(res *pb.ApiResponsesMerchantPolicies) *model.APIResponsesMerchantPolicy
	ToGraphqlResponsePaginationMerchantPolicyDeleteAt(res *pb.ApiResponsePaginationMerchantPoliciesDeleteAt) *model.APIResponsePaginationMerchantPolicyDeleteAt
	ToGraphqlResponsePaginationMerchantPolicy(res *pb.ApiResponsePaginationMerchantPolicies) *model.APIResponsePaginationMerchantPolicy
}
