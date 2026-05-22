package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type MerchantBusinessQueryHandler interface {
	pb.MerchantBusinessQueryServiceServer
}

type MerchantBusinessCommandHandler interface {
	pb.MerchantBusinessCommandServiceServer
}
