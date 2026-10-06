package user_test

import (
	"fmt"
	"net/http"
	"testing"

	db "github.com/MamangRust/monolith-graphql-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
	user_cache "github.com/MamangRust/monolith-graphql-ecommerce-user/cache"
	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	gapi "github.com/MamangRust/monolith-graphql-ecommerce-user/handler"
	"github.com/MamangRust/monolith-graphql-ecommerce-user/repository"
	"github.com/MamangRust/monolith-graphql-ecommerce-user/service"

	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"

	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
	pbuserrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/user_role"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
)

type UserHandlerTestSuite struct {
	tests.BaseTestSuite
	client    pbuser.UserCommandServiceClient
	gql       http.Handler
	userID    int
	userEmail string
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()
	roleClient := pbrole.NewRoleQueryServiceClient(s.Conns["role"])

	queries := db.New(s.DBPool())
	repos := repository.NewRepositories(&repository.Deps{Db: queries, RoleQueryClient: roleClient, UserRoleClient: pbuserrole.NewUserRoleServiceClient(s.Conns["role"])})

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheStore := s.GetCacheStore()
	mencache := user_cache.NewMencache(cacheStore)

	userService := service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Hash:          hasher,
		Logger:        log,
		Observability: s.Obs,
	})

	// Start gRPC Server
	userHandler := gapi.NewHandler(&gapi.Deps{
		Service: userService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pbuser.RegisterUserQueryServiceServer(server, userHandler.UserQuery)
	pbuser.RegisterUserCommandServiceServer(server, userHandler.UserCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)
	s.Conns["user"] = conn
	s.client = pbuser.NewUserCommandServiceClient(conn)

	s.gql = tests.NewGraphQLHandler(s.Conns, cacheStore, log)
}

func (s *UserHandlerTestSuite) TestUserApiLifecycle() {
	// 1. Create
	s.userEmail = "handler.user@example.com"
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createUser(input: {
			firstname: "Handler"
			lastname: "User"
			email: "%s"
			password: "password123"
			confirm_password: "password123"
		}) {
			status message
			data { id email firstname }
		}
	}`, s.userEmail))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createUser")
	s.Require().NotNil(data)
	s.userID = int(data["id"].(float64))
	s.Require().NotZero(s.userID)
	s.Equal(s.userEmail, data["email"])

	// 2. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllUsers(input: { page: 1, page_size: 10 }) { status message data { id email } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllUsers"))

	// 3. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findByIdUser(input: { id: %d }) { status message data { id email } } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findByIdUser")
	s.Require().NotNil(data)
	s.Equal(float64(s.userID), data["id"])

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findByActiveUsers(input: { page: 1, page_size: 10 }) { status message data { id email } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByActiveUsers"))

	// 5. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateUser(input: {
			id: %d
			firstname: "Updated"
			lastname: "User"
			email: "%s"
			password: "password123"
			confirm_password: "password123"
		}) {
			status message
			data { id firstname email }
		}
	}`, s.userID, s.userEmail))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateUser")
	s.Require().NotNil(data)
	s.Equal("Updated", data["firstname"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedUser(input: { id: %d }) { status message data { id } } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findByTrashedUsers(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByTrashedUsers"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreUser(input: { id: %d }) { status message data { id } } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedUser(input: { id: %d }) { status message } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteUserPermanent(input: { id: %d }) { status message } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllUser { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllUserPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestUserHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserHandlerTestSuite))
}
