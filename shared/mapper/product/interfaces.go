package productapimapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/response"

	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
)

type ProductBaseResponseMapper interface {
	ToResponseProduct(product *pbproduct.ProductResponse) *response.ProductResponse
	ToResponsesProduct(products []*pbproduct.ProductResponse) []*response.ProductResponse
	ToResponseProductDeleteAt(product *pbproduct.ProductResponseDeleteAt) *response.ProductResponseDeleteAt
	ToResponsesProductDeleteAt(products []*pbproduct.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt
	ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct
	ToApiResponsePaginationProductDeleteAt(pbResponse *pbproduct.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt
}

type ProductQueryResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProduct(pbResponse *pbproduct.ApiResponsesProduct) *response.ApiResponsesProduct
	ToApiResponsePaginationProduct(pbResponse *pbproduct.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct
}

type ProductCommandResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProductDeleteAt(pbResponse *pbproduct.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt
	ToApiResponseProductDelete(pbResponse *pbproduct.ApiResponseProductDelete) *response.ApiResponseProductDelete
	ToApiResponseProductAll(pbResponse *pbproduct.ApiResponseProductAll) *response.ApiResponseProductAll
}
