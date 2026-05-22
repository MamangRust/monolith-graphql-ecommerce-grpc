package cartgraphqlmapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type CartGraphqlMapper interface {
	ToGraphqlResponseCartDelete(res *pb.ApiResponseCartDelete) *model.APIResponseCartDelete
	ToGraphqlResponseCartAll(res *pb.ApiResponseCartAll) *model.APIResponseCartAll
	ToGraphqlResponseCart(res *pb.ApiResponseCart) *model.APIResponseCart
	ToGraphqlResponsePaginationCart(res *pb.ApiResponsePaginationCart) *model.APIResponsePaginationCart
}
