package handler

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-user/service"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	UserQuery   UserQueryHandler
	UserCommand UserCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		UserQuery:   NewUserQueryHandler(deps.Service.UserQuery, deps.Logger),
		UserCommand: NewUserCommandHandler(deps.Service.UserCommand, deps.Logger),
	}
}
