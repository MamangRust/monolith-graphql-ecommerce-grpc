package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	pb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	"github.com/MamangRust/monolith-graphql-ecommerce-role/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	RoleQuery   pb.RoleQueryServiceServer
	RoleCommand pb.RoleCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		RoleQuery:   NewRoleQueryHandler(deps.Service.RoleQuery, deps.Logger),
		RoleCommand: NewRoleCommandHandler(deps.Service.RoleCommand, deps.Logger),
	}
}
