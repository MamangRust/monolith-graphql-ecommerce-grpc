package reviewdetailapimapper

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-graphql-ecommerce-shared/mapper/pagination"

	pbreview_detail "github.com/MamangRust/monolith-graphql-ecommerce-pb/review_detail"
)

type reviewDetailCommandResponseMapper struct{}

func NewReviewDetailCommandResponseMapper() ReviewDetailCommandResponseMapper {
	return &reviewDetailCommandResponseMapper{}
}

func (m *reviewDetailCommandResponseMapper) ToResponseReviewDetail(reviewDetail *pbreview_detail.ReviewDetailsResponse) *response.ReviewDetailsResponse {
	return &response.ReviewDetailsResponse{
		ID:        int(reviewDetail.Id),
		ReviewID:  int(reviewDetail.ReviewId),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
	}
}

func (m *reviewDetailCommandResponseMapper) ToResponsesReviewDetail(ReviewDetails []*pbreview_detail.ReviewDetailsResponse) []*response.ReviewDetailsResponse {
	var mappedReviewDetails []*response.ReviewDetailsResponse
	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetail(ReviewDetail))
	}
	return mappedReviewDetails
}

func (m *reviewDetailCommandResponseMapper) ToResponseReviewDetailDeleteAt(reviewDetail *pbreview_detail.ReviewDetailsResponseDeleteAt) *response.ReviewDetailsResponseDeleteAt {
	var deletedAt string
	if reviewDetail.DeletedAt != nil {
		deletedAt = reviewDetail.DeletedAt.Value
	}

	return &response.ReviewDetailsResponseDeleteAt{
		ID:        int(reviewDetail.Id),
		ReviewID:  int(reviewDetail.ReviewId),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (m *reviewDetailCommandResponseMapper) ToResponsesReviewDetailDeleteAt(ReviewDetails []*pbreview_detail.ReviewDetailsResponseDeleteAt) []*response.ReviewDetailsResponseDeleteAt {
	var mappedReviewDetails []*response.ReviewDetailsResponseDeleteAt
	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetailDeleteAt(ReviewDetail))
	}
	return mappedReviewDetails
}

func (m *reviewDetailCommandResponseMapper) ToApiResponseReviewDetailDeleteAt(pbResponse *pbreview_detail.ApiResponseReviewDetailDeleteAt) *response.ApiResponseReviewDetailDeleteAt {
	return &response.ApiResponseReviewDetailDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseReviewDetailDeleteAt(pbResponse.Data),
	}
}

func (m *reviewDetailCommandResponseMapper) ToApiResponsePaginationReviewDetailDeleteAt(pbResponse *pbreview_detail.ApiResponsePaginationReviewDetailsDeleteAt) *response.ApiResponsePaginationReviewDetailsDeleteAt {
	return &response.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesReviewDetailDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
