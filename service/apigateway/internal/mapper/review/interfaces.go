package reviewgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/review"
)

type ReviewGraphqlMapper interface {
	ToGraphqlResponseReview(res *pb.ApiResponseReview) *model.APIResponseReview
	ToGraphqlResponseReviewDeleteAt(res *pb.ApiResponseReviewDeleteAt) *model.APIResponseReviewDeleteAt
	ToGraphqlResponsesReview(res *pb.ApiResponsesReview) *model.APIResponsesReview
	ToGraphqlResponseReviewDelete(res *pb.ApiResponseReviewDelete) *model.APIResponseReviewDelete
	ToGraphqlResponseReviewAll(res *pb.ApiResponseReviewAll) *model.APIResponseReviewAll
	ToGraphqlResponsePaginationReviewDeleteAt(res *pb.ApiResponsePaginationReviewDeleteAt) *model.APIResponsePaginationReviewDeleteAt
	ToGraphqlResponsePaginationReview(res *pb.ApiResponsePaginationReview) *model.APIResponsePaginationReview
	ToGraphqlResponsePaginationReviewRelationDetail(res *pb.ApiResponsePaginationReviewDetail) *model.APIResponsePaginationReviewRelationDetail
}
