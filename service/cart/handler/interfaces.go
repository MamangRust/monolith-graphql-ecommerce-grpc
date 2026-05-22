package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type CartQueryHandler interface {
	pb.CartQueryServiceServer
}

type CartCommandHandler interface {
	pb.CartCommandServiceServer
}
