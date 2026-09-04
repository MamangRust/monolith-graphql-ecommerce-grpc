package review_detailgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pbreview "github.com/MamangRust/monolith-graphql-ecommerce-pb/review"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/review_detail"
)

type ReviewDetailGraphqlMapper interface {
	ToGraphqlResponseDelete(res *pbreview.ApiResponseReviewDelete) *model.APIResponseReviewDetailDelete
	ToGraphqlResponseAll(res *pbreview.ApiResponseReviewAll) *model.APIResponseReviewDetailAll
	ToGraphqlResponseReviewDetail(res *pb.ApiResponseReviewDetail) *model.APIResponseReviewDetail
	ToGraphqlResponsesReviewDetail(res *pb.ApiResponsesReviewDetails) *model.APIResponsesReviewDetails
	ToGraphqlResponseReviewDetailDeleteAt(res *pb.ApiResponseReviewDetailDeleteAt) *model.APIResponseReviewDetailDeleteAt
	ToGraphqlResponsePaginationReviewDetail(res *pb.ApiResponsePaginationReviewDetails) *model.APIResponsePaginationReviewDetails
	ToGraphqlResponsePaginationReviewDetailDeleteAt(res *pb.ApiResponsePaginationReviewDetailsDeleteAt) *model.APIResponsePaginationReviewDetailsDeleteAt
}
