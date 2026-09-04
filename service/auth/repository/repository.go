package repository

import (
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

func NewRepositories(DB *db.Queries,
	userQuery pbuser.UserQueryServiceClient,
	userCommand pbuser.UserCommandServiceClient,
	roleQuery pbrole.RoleQueryServiceClient,
	roleCommand pbrole.RoleCommandServiceClient,
) *Repositories {
	return &Repositories{
		User:         NewUserRepository(userQuery, userCommand),
		RefreshToken: NewRefreshTokenRepository(DB),
		UserRole:     NewUserRoleRepository(roleCommand),
		Role:         NewRoleRepository(roleQuery),
		ResetToken:   NewResetTokenRepository(DB),
	}
}
