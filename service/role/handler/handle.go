package handler

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-role/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
	pbuserrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/user_role"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	RoleQuery   pbrole.RoleQueryServiceServer
	RoleCommand pbrole.RoleCommandServiceServer
	UserRole    pbuserrole.UserRoleServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		RoleQuery:   NewRoleQueryHandler(deps.Service.RoleQuery, deps.Logger),
		RoleCommand: NewRoleCommandHandler(deps.Service.RoleCommand, deps.Logger),
		UserRole:    NewUserRoleHandler(deps.Service, deps.Logger),
	}
}
