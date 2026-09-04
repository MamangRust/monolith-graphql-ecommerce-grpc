package paginationapimapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/response"

	pbcommon "github.com/MamangRust/monolith-graphql-ecommerce-pb/common"
)

func MapPaginationMeta(s *pbcommon.PaginationMeta) *response.PaginationMeta {
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
