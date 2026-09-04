package handler

import (

	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type UserQueryHandler interface {
	pbuser.UserQueryServiceServer
}

type UserCommandHandler interface {
	pbuser.UserCommandServiceServer
}
