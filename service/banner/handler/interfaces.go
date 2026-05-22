package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type BannerQueryHandler interface {
	pb.BannerQueryServiceServer
}

type BannerCommandHandler interface {
	pb.BannerCommandServiceServer
}
