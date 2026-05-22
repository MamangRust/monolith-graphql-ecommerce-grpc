package paginationapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

func MapPaginationMeta(s *pb.PaginationMeta) *response.PaginationMeta {
	if s == nil {
		return nil
	}
	return &response.PaginationMeta{
		CurrentPage:  int(s.CurrentPage),
		PageSize:     int(s.PageSize),
		TotalRecords: int(s.TotalRecords),
		TotalPages:   int(s.TotalPages),
	}
}
