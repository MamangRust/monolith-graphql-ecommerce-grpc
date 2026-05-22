package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type ProductQueryHandler interface {
	pb.ProductQueryServiceServer
}

type ProductCommandHandler interface {
	pb.ProductCommandServiceServer
}
