package handler

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-role/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	RoleQuery   pbrole.RoleQueryServiceServer
	RoleCommand pbrole.RoleCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		RoleQuery:   NewRoleQueryHandler(deps.Service.RoleQuery, deps.Logger),
		RoleCommand: NewRoleCommandHandler(deps.Service.RoleCommand, deps.Logger),
	}
}
