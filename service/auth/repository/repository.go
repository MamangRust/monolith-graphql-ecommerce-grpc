package repository

import (
	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
	pbuserrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/user_role"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"

	roleadapter "github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter/role"
	useradapter "github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/monolith-graphql-ecommerce-pkg/adapter/user_role"
)

// Repositories bundles all auth repositories. It uses named fields (not
// embedding) because UserRepository and RoleRepository both declare FindById,
// and RefreshTokenRepository and ResetTokenRepository both declare FindByToken.
type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

// GuardOptions collects the dependency guards for each remote dependency. Role
// reads and user_role writes get separate breakers so one failure mode cannot
// starve the other.
type GuardOptions struct {
	User     []adapter.GuardOption
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

// Compile-time assertions: the shared adapters satisfy auth's repository
// contracts directly.
var (
	_ UserRepository     = (*useradapter.Repository)(nil)
	_ RoleRepository     = (*roleadapter.Repository)(nil)
	_ UserRoleRepository = (*userroleadapter.Repository)(nil)
)

type Deps struct {
	Db                *db.Queries
	UserQueryClient   pbuser.UserQueryServiceClient
	UserCommandClient pbuser.UserCommandServiceClient
	RoleQueryClient   pbrole.RoleQueryServiceClient
	RoleCommandClient pbrole.RoleCommandServiceClient
	UserRoleClient    pbuserrole.UserRoleServiceClient
	Guard             GuardOptions
}

func NewRepositories(
	deps *Deps,
) *Repositories {
	return &Repositories{
		User:         useradapter.New(deps.UserQueryClient, deps.UserCommandClient, deps.Guard.User...),
		RefreshToken: NewRefreshTokenRepository(deps.Db),
		UserRole:     userroleadapter.New(deps.UserRoleClient, deps.Guard.UserRole...),
		Role:         roleadapter.New(deps.RoleQueryClient, deps.Guard.Role...),
		ResetToken:   NewResetTokenRepository(deps.Db),
	}
}
