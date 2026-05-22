package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleRepository
}

func NewRepositories(DB *db.Queries, roleClient pb.RoleQueryServiceClient) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(DB),
		UserQuery:   NewUserQueryRepository(DB),
		Role:        NewRoleRepository(roleClient),
	}
}
