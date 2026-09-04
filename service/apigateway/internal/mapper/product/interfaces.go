package productgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
)

type ProductGraphqlMapper interface {
	ToGraphqlResponseProduct(res *pb.ApiResponseProduct) *model.APIResponseProduct
	ToGraphqlResponsesProduct(res *pb.ApiResponsesProduct) *model.APIResponsesProduct
	ToGraphqlResponseProductDeleteAt(res *pb.ApiResponseProductDeleteAt) *model.APIResponseProductDeleteAt
	ToGraphqlResponseProductDelete(res *pb.ApiResponseProductDelete) *model.APIResponseProductDelete
	ToGraphqlResponseProductAll(res *pb.ApiResponseProductAll) *model.APIResponseProductAll
	ToGraphqlResponsePaginationProduct(res *pb.ApiResponsePaginationProduct) *model.APIResponsePaginationProduct
	ToGraphqlResponsePaginationProductDeleteAt(res *pb.ApiResponsePaginationProductDeleteAt) *model.APIResponsePaginationProductDeleteAt
}
