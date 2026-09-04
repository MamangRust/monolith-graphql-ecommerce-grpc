package role_test

import (
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type RoleApiTestSuite struct {
	tests.BaseTestSuite
	gql    http.Handler
	roleID int
}

func (s *RoleApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *RoleApiTestSuite) TestRoleApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, `mutation {
		createRole(input: { name: "API Role" }) {
			status message
			data { id name }
		}
	}`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createRole")
	s.Require().NotNil(data)
	s.Equal("API Role", data["name"])
	s.roleID = int(data["id"].(float64))
	s.Require().NotZero(s.roleID)

	// 2. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllRole(input: { page: 1, page_size: 10 }) { status message data { id name } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllRole"))

	// 3. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findByIdRole(input: { role_id: %d }) { status message data { id name } } }`, s.roleID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findByIdRole")
	s.Require().NotNil(data)
	s.Equal(float64(s.roleID), data["id"])

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findByActiveRole(input: { page: 1, page_size: 10 }) { status message data { id name } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByActiveRole"))

	// 5. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findByTrashedRole(input: { page: 1, page_size: 10 }) { status message data { id name } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByTrashedRole"))

	// 6. FindByUserId
	body, err = tests.DoGraphQL(s.gql, `query { findByUserIdRole(input: { user_id: 1 }) { status message data { id name } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByUserIdRole"))

	// 7. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateRole(input: { id: %d, name: "Updated API Role" }) {
			status message
			data { id name }
		}
	}`, s.roleID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateRole")
	s.Require().NotNil(data)
	s.Equal("Updated API Role", data["name"])

	// 8. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedRole(input: { role_id: %d }) { status message data { id } } }`, s.roleID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreRole(input: { role_id: %d }) { status message data { id } } }`, s.roleID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedRole(input: { role_id: %d }) { status message } }`, s.roleID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteRolePermanent(input: { role_id: %d }) { status message } }`, s.roleID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllRole { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 12. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllRolePermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestRoleApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleApiTestSuite))
}
