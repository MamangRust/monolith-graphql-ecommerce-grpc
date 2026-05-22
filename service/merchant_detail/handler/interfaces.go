package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type MerchantDetailQueryHandler interface {
	pb.MerchantDetailQueryServiceServer
}

type MerchantDetailCommandHandler interface {
	pb.MerchantDetailCommandServiceServer
}

type MerchantSocialLinkCommandHandler interface {
	pb.MerchantSocialLinkServiceServer
}
