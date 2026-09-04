package cartapimapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/response"

	pbcart "github.com/MamangRust/monolith-graphql-ecommerce-pb/cart"
)

type CartBaseResponseMapper interface {
	ToResponseCart(pbResponse *pbcart.CartResponse) *response.CartResponse
	ToResponseCarts(pbResponse []*pbcart.CartResponse) []*response.CartResponse
	ToApiResponseCart(pbResponse *pbcart.ApiResponseCart) *response.ApiResponseCart
}

type CartQueryResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartPagination(pbResponse *pbcart.ApiResponsePaginationCart) *response.ApiResponseCartPagination
}

type CartCommandResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartDelete(pbResponse *pbcart.ApiResponseCartDelete) *response.ApiResponseCartDelete
	ToApiResponseCartAll(pbResponse *pbcart.ApiResponseCartAll) *response.ApiResponseCartAll
}
