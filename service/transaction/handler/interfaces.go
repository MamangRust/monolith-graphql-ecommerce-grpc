package handler

import pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"

type TransactionQueryHandler interface {
	pb.TransactionQueryServiceServer
}

type TransactionCommandHandler interface {
	pb.TransactionCommandServiceServer
}

type TransactionStatsHandler interface {
	pb.TransactionStatsServiceServer
}

type TransactionStatsByMerchantHandler interface {
	pb.TransactionStatsByMerchantServiceServer
}
