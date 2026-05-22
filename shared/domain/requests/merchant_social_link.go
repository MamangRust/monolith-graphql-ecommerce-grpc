package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantSocialRequest struct {
	Platform string `json:"platform" validate:"required"`
	Url      string `json:"url" validate:"required,url"`
}

type UpdateMerchantSocialRequest struct {
	ID               int    `json:"id" validate:"required"`
	MerchantDetailID *int   `json:"merchant_detail_id" validate:"required"`
	Platform         string `json:"platform" validate:"required"`
	Url              string `json:"url" validate:"required,url"`
}

type CreateBatchMerchantSocialRequest struct {
	MerchantDetailID int                            `json:"merchant_detail_id" validate:"required"`
	SocialLinks      []*CreateMerchantSocialRequest `json:"social_links"`
}

type UpdateBatchMerchantSocialRequest struct {
	MerchantDetailID int                            `json:"merchant_detail_id" validate:"required"`
	SocialLinks      []*UpdateMerchantSocialRequest `json:"social_links"`
}

func (r *CreateBatchMerchantSocialRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateBatchMerchantSocialRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
