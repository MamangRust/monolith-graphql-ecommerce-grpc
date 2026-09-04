package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

func NewRepositories(DB *db.Queries, userQuery pbuser.UserQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(DB),
		MerchantCommand:         NewMerchantCommandRepository(DB),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(DB),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(DB),
		UserQuery:               NewUserQueryRepository(userQuery),
	}
}
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
