package handler

import (
	pbbanner "github.com/MamangRust/monolith-graphql-ecommerce-pb/banner"
)
type BannerQueryHandler interface {
	pbbanner.BannerQueryServiceServer
}

type BannerCommandHandler interface {
	pbbanner.BannerCommandServiceServer
}

