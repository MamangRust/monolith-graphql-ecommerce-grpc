package merchant_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	merchantID int
	userID     int
}

func (s *MerchantApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupMerchantService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *MerchantApiTestSuite) TestMerchantApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createMerchant(input: {
			user_id: %d
			name: "Test Merchant"
			description: "Test Description"
			address: "Test Address"
			contact_email: "merchant@example.com"
			contact_phone: "123456789"
			status: "active"
		}) {
			status message
			data { id name }
		}
	}`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createMerchant")
	s.Require().NotNil(data)
	s.Equal("Test Merchant", data["name"])
	s.merchantID = int(data["id"].(float64))
	s.Require().NotZero(s.merchantID)

	// 2. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findMerchantById(input: { id: %d }) { status message data { id name } } }`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findMerchantById")
	s.Require().NotNil(data)
	s.Equal(float64(s.merchantID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllMerchants(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllMerchants"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveMerchants(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveMerchants"))

	// 5. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateMerchant(input: {
			merchant_id: %d
			user_id: %d
			name: "Updated Merchant"
			description: "Updated Description"
			address: "Updated Address"
			contact_email: "updated@example.com"
			contact_phone: "987654321"
			status: "active"
		}) {
			status message
			data { id name }
		}
	}`, s.merchantID, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateMerchant")
	s.Require().NotNil(data)
	s.Equal("Updated Merchant", data["name"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchant(input: { id: %d }) { status message data { id } } }`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedMerchants(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedMerchants"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreMerchant(input: { id: %d }) { status message data { id } } }`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchant(input: { id: %d }) { status message } }`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteMerchantPermanent(input: { id: %d }) { status message } }`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllMerchants { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllMerchantsPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestMerchantApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantApiTestSuite))
}
