package pagination

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
	pbcommon "github.com/MamangRust/monolith-graphql-ecommerce-pb/common"
)

func MapPaginationMeta(meta *pbcommon.PaginationMeta) *model.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &model.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		TotalPages:   meta.TotalPages,
		PageSize:     meta.PageSize,
		TotalRecords: meta.TotalRecords,
	}
}
