package repository

import (
	"context"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	user_errors "github.com/MamangRust/monolith-graphql-ecommerce-shared/errors/user_errors"

	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type userQueryRepository struct {
	client pbuser.UserQueryServiceClient
}

func NewUserQueryRepository(client pbuser.UserQueryServiceClient) *userQueryRepository {
	return &userQueryRepository{
		client: client,
	}
}

func (r *userQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.client.FindById(ctx, &pbuser.FindByIdUserRequest{Id: int32(user_id)})
	if err != nil {
		return nil, user_errors.ErrUserInternal.WithInternal(err)
	}

	return &db.GetUserByIDRow{
		UserID:    res.Data.Id,
		Firstname: "", // UserResponse does not provide Firstname
		Lastname:  "", // UserResponse does not provide Lastname
		Email:     res.Data.Email,
	}, nil
}
