package auth_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-graphql-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-auth/service"
	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	mencache "github.com/MamangRust/monolith-graphql-ecommerce-auth/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"

	user_handler "github.com/MamangRust/monolith-graphql-ecommerce-user/handler"
	user_service "github.com/MamangRust/monolith-graphql-ecommerce-user/service"
	user_repo "github.com/MamangRust/monolith-graphql-ecommerce-user/repository"
	role_handler "github.com/MamangRust/monolith-graphql-ecommerce-role/handler"
	role_service "github.com/MamangRust/monolith-graphql-ecommerce-role/service"
	role_repo "github.com/MamangRust/monolith-graphql-ecommerce-role/repository"
	user_cache "github.com/MamangRust/monolith-graphql-ecommerce-user/cache"
	role_cache "github.com/MamangRust/monolith-graphql-ecommerce-role/cache"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
	pbuserrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/user_role"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	service     *service.Service
	email       string
	password    string
}

func (s *AuthServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	// 1. Setup Role Service & gRPC Server
	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(queries)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repository:    roleRepos,
		Logger:        log,
		Cache:         roleMencache,
		Observability: obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  log,
	})
	roleServer := grpc.NewServer()
	pbrole.RegisterRoleQueryServiceServer(roleServer, roleGapi.RoleQuery)
	pbrole.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	pbuserrole.RegisterUserRoleServiceServer(roleServer, roleGapi.UserRole)
	roleLis, _ := net.Listen("tcp", "localhost:0")
	go roleServer.Serve(roleLis)
	roleConn, _ := grpc.NewClient(roleLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 2. Setup User Service & gRPC Server
	userMencache := user_cache.NewMencache(cacheStore)
	roleQueryClientForUser := pbrole.NewRoleQueryServiceClient(roleConn)
	userRepos := user_repo.NewRepositories(&user_repo.Deps{Db: queries, RoleQueryClient: roleQueryClientForUser, UserRoleClient: pbuserrole.NewUserRoleServiceClient(roleConn)})
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        log,
		Hash:         hasher,
		Cache:         userMencache,
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	pbuser.RegisterUserQueryServiceServer(userServer, userGapi.UserQuery)
	pbuser.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)
	userLis, _ := net.Listen("tcp", "localhost:0")
	go userServer.Serve(userLis)
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 3. Setup Auth Service with gRPC clients
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)
	userCommandClient := pbuser.NewUserCommandServiceClient(userConn)
	roleQueryClient := pbrole.NewRoleQueryServiceClient(roleConn)
	roleCommandClient := pbrole.NewRoleCommandServiceClient(roleConn)

	repos := repository.NewRepositories(&repository.Deps{Db: queries, UserQueryClient: userQueryClient, UserCommandClient: userCommandClient, RoleQueryClient: roleQueryClient, RoleCommandClient: roleCommandClient, UserRoleClient: pbuserrole.NewUserRoleServiceClient(roleConn)})

	tokenManager, _ := auth.NewManager("mysecret")
	s.service = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache.NewMencache(cacheStore),
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	s.email = "auth.service.test@example.com"
	s.password = "password123"

	// Seed ROLE_ADMIN
	_, _ = pool.Exec(context.Background(), "INSERT INTO roles (role_name) VALUES ('ROLE_ADMIN')")
}

func (s *AuthServiceTestSuite) TearDownSuite() {
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *AuthServiceTestSuite) Test1_Register() {
	ctx := context.Background()
	req := &requests.RegisterRequest{
		FirstName:       "Auth",
		LastName:        "Service",
		Email:           s.email,
		Password:        s.password,
		ConfirmPassword: s.password,
	}

	res, err := s.service.Register.Register(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(s.email, res.Email)
}

func (s *AuthServiceTestSuite) Test2_Login() {
	ctx := context.Background()
	req := &requests.AuthRequest{
		Email:    s.email,
		Password: s.password,
	}

	res, err := s.service.Login.Login(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.NotEmpty(res.AccessToken)
	s.NotEmpty(res.RefreshToken)
}

func (s *AuthServiceTestSuite) Test4_LoginLockout() {
	ctx := context.Background()
	email := "locked.user@example.com"
	password := "wrongpassword"

	// Register user first
	regReq := &requests.RegisterRequest{
		FirstName:       "Locked",
		LastName:        "User",
		Email:           email,
		Password:        "correctpassword",
		ConfirmPassword: "correctpassword",
	}
	_, err := s.service.Register.Register(ctx, regReq)
	s.NoError(err)

	loginReq := &requests.AuthRequest{
		Email:    email,
		Password: password,
	}

	// Fail login 5 times
	for i := 0; i < 5; i++ {
		_, err := s.service.Login.Login(ctx, loginReq)
		s.Error(err)
	}

	// 6th attempt should return ErrAccountLocked
	_, err = s.service.Login.Login(ctx, loginReq)
	s.Error(err)
	s.Contains(err.Error(), "Account temporarily locked")
}

func (s *AuthServiceTestSuite) Test3_ForgotPassword() {
	ctx := context.Background()
	
	success, err := s.service.PasswordReset.ForgotPassword(ctx, s.email)
	s.NoError(err)
	s.True(success)
}

func TestAuthServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthServiceTestSuite))
}
