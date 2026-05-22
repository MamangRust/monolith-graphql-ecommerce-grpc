package handler

import (
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type UserQueryHandler interface {
	pb.UserQueryServiceServer
}

type UserCommandHandler interface {
	pb.UserCommandServiceServer
}
