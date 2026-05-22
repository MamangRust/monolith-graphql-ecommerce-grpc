package productapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type ProductBaseResponseMapper interface {
	ToResponseProduct(product *pb.ProductResponse) *response.ProductResponse
	ToResponsesProduct(products []*pb.ProductResponse) []*response.ProductResponse
	ToResponseProductDeleteAt(product *pb.ProductResponseDeleteAt) *response.ProductResponseDeleteAt
	ToResponsesProductDeleteAt(products []*pb.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt
	ToApiResponseProduct(pbResponse *pb.ApiResponseProduct) *response.ApiResponseProduct
	ToApiResponsePaginationProductDeleteAt(pbResponse *pb.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt
}

type ProductQueryResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProduct(pbResponse *pb.ApiResponsesProduct) *response.ApiResponsesProduct
	ToApiResponsePaginationProduct(pbResponse *pb.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct
}

type ProductCommandResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProductDeleteAt(pbResponse *pb.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt
	ToApiResponseProductDelete(pbResponse *pb.ApiResponseProductDelete) *response.ApiResponseProductDelete
	ToApiResponseProductAll(pbResponse *pb.ApiResponseProductAll) *response.ApiResponseProductAll
}
