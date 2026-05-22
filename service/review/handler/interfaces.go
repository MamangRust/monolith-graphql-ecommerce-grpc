package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type ReviewHandleGrpc interface {
	pb.ReviewQueryServiceServer
	pb.ReviewCommandServiceServer
}

type ReviewQueryHandler interface {
	pb.ReviewQueryServiceServer
}

type ReviewCommandHandler interface {
	pb.ReviewCommandServiceServer
}
