package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type MerchantPolicyQueryHandler interface {
	pb.MerchantPolicyQueryServiceServer
}

type MerchantPolicyCommandHandler interface {
	pb.MerchantPolicyCommandServiceServer
}
