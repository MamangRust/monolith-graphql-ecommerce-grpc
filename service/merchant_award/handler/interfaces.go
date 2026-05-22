package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type MerchantAwardQueryHandler interface {
	pb.MerchantAwardQueryServiceServer
}

type MerchantAwardCommandHandler interface {
	pb.MerchantAwardCommandServiceServer
}
