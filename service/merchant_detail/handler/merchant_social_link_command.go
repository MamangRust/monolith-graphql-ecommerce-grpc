package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantsociallink_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	"github.com/MamangRust/monolith-graphql-ecommerce-merchant_detail/service"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type merchantSocialLinkCommandHandler struct {
	pb.UnimplementedMerchantSocialLinkServiceServer
	MerchantSocialLinkCommand service.MerchantSocialLinkCommandService
	MerchantDetailQuery       service.MerchantDetailQueryService
	logger                    logger.LoggerInterface
}

func NewMerchantSocialLinkCommandHandler(
	svc service.MerchantSocialLinkCommandService,
	querySvc service.MerchantDetailQueryService,
	logger logger.LoggerInterface,
) MerchantSocialLinkCommandHandler {
	return &merchantSocialLinkCommandHandler{
		MerchantSocialLinkCommand: svc,
		MerchantDetailQuery:       querySvc,
		logger:                    logger,
	}
}

func (s *merchantSocialLinkCommandHandler) CreateMerchantSocialLink(ctx context.Context, request *pb.CreateMerchantSocialInput) (*pb.ApiResponseMerchantSocialMediaLink, error) {
	if request.MerchantDetailId <= 0 {
		return nil, merchantsociallink_errors.ErrGrpcMerchantSocialLinkInvalidId
	}
	if len(request.SocialLinks) == 0 {
		return nil, merchantsociallink_errors.ErrGrpcValidateCreateMerchantSocialLink
	}

	socialLinks := make([]*requests.CreateMerchantSocialRequest, 0, len(request.SocialLinks))
	for _, link := range request.SocialLinks {
		socialLinks = append(socialLinks, &requests.CreateMerchantSocialRequest{
			Platform: link.Platform,
			Url:      link.Url,
		})
	}

	req := &requests.CreateBatchMerchantSocialRequest{
		MerchantDetailID: int(request.MerchantDetailId),
		SocialLinks:      socialLinks,
	}

	if err := req.Validate(); err != nil {
		return nil, merchantsociallink_errors.ErrGrpcValidateCreateMerchantSocialLink
	}

	createdLinks, err := s.MerchantSocialLinkCommand.CreateSocialLink(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoLinks := make([]*pb.MerchantSocialMediaLinkResponse, 0, len(createdLinks))
	for _, link := range createdLinks {
		protoLinks = append(protoLinks, &pb.MerchantSocialMediaLinkResponse{
			Id:       int32(link.MerchantSocialID),
			Platform: link.Platform,
			Url:      link.Url,
		})
	}

	return &pb.ApiResponseMerchantSocialMediaLink{
		Status:  "success",
		Message: "Successfully created merchant social links",
		Data:    protoLinks,
	}, nil
}

func (s *merchantSocialLinkCommandHandler) UpdateMerchantSocialLink(ctx context.Context, request *pb.UpdateMerchantSocialInput) (*pb.ApiResponseMerchantSocialMediaLink, error) {
	if request.MerchantDetailId <= 0 {
		return nil, merchantsociallink_errors.ErrGrpcMerchantSocialLinkInvalidId
	}
	if len(request.SocialLinks) == 0 {
		return nil, merchantsociallink_errors.ErrGrpcValidateUpdateMerchantSocialLink
	}

	socialLinks := make([]*requests.UpdateMerchantSocialRequest, 0, len(request.SocialLinks))
	for _, link := range request.SocialLinks {
		socialLinks = append(socialLinks, &requests.UpdateMerchantSocialRequest{
			Platform: link.Platform,
			Url:      link.Url,
		})
	}

	req := &requests.UpdateBatchMerchantSocialRequest{
		MerchantDetailID: int(request.MerchantDetailId),
		SocialLinks:      socialLinks,
	}

	if err := req.Validate(); err != nil {
		return nil, merchantsociallink_errors.ErrGrpcValidateUpdateMerchantSocialLink
	}

	updatedLinks, err := s.MerchantSocialLinkCommand.UpdateSocialLink(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoLinks := make([]*pb.MerchantSocialMediaLinkResponse, 0, len(updatedLinks))
	for _, link := range updatedLinks {
		protoLinks = append(protoLinks, &pb.MerchantSocialMediaLinkResponse{
			Id:       int32(link.MerchantSocialID),
			Platform: link.Platform,
			Url:      link.Url,
		})
	}

	return &pb.ApiResponseMerchantSocialMediaLink{
		Status:  "success",
		Message: "Successfully updated merchant social links",
		Data:    protoLinks,
	}, nil
}
