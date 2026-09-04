package shipping_addressgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
)

type ShippingAddresGraphqlMapper interface {
	ToGraphResponseShippingAddress(res *pb.ApiResponseShipping) *model.APIResponseShipping
	ToGraphResponseShippingAddressDeleteAt(res *pb.ApiResponseShippingDeleteAt) *model.APIResponseShippingDeleteAt
	ToGraphResponsesShippingAddress(res *pb.ApiResponsesShipping) *model.APIResponsesShipping
	ToGraphResponseShippingAddressDelete(res *pb.ApiResponseShippingDelete) *model.APIResponseShippingDelete
	ToGraphResponseShippingAddressAll(res *pb.ApiResponseShippingAll) *model.APIResponseShippingAll
	ToGraphResponsePaginationShippingAddressDeleteAt(res *pb.ApiResponsePaginationShippingDeleteAt) *model.APIResponsePaginationShippingDeleteAt
	ToGraphResponsePaginationShippingAddress(res *pb.ApiResponsePaginationShipping) *model.APIResponsePaginationShipping
}
