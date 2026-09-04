package repository

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
)

type roleRepository struct {
	client pbrole.RoleQueryServiceClient
}

func NewRoleRepository(client pbrole.RoleQueryServiceClient) RoleRepository {
	return &roleRepository{
		client: client,
	}
}

func (r *roleRepository) FindById(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.client.FindByIdRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: int32(id)})
	if err != nil {
		return nil, fmt.Errorf("failed to find role by ID %d: %w", id, err)
	}

	return r.mapToDBRole(res.Data), nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.client.FindByNameRole(ctx, &pbrole.FindByNameRoleRequest{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find role by name %s: %w", name, err)
	}

	return r.mapToDBRole(res.Data), nil
}

func (r *roleRepository) mapToDBRole(pbRole *pbrole.RoleResponse) *db.Role {
	if pbRole == nil {
		return nil
	}

	createdAt, _ := time.Parse(time.RFC3339, pbRole.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbRole.UpdatedAt)

	return &db.Role{
		RoleID:   pbRole.Id,
		RoleName: pbRole.Name,
		CreatedAt: pgtype.Timestamp{
			Time:  createdAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  updatedAt,
			Valid: true,
		},
	}
}

