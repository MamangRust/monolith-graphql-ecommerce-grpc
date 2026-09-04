package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleRepository
}

func NewRepositories(DB *db.Queries, roleClient pbrole.RoleQueryServiceClient) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(DB),
		UserQuery:   NewUserQueryRepository(DB),
		Role:        NewRoleRepository(roleClient),
	}
}
