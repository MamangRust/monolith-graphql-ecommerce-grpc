package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type OrderItemQueryHandler interface {
	pb.OrderItemQueryServiceServer
}

type OrderItemCommandHandler interface {
	pb.OrderItemCommandServiceServer
}
