package merchant_business_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantBusinessApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	businessID int
	merchantID int
	userID     int
}

func (s *MerchantBusinessApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupMerchantBusinessService()

	// Seed dependencies
	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *MerchantBusinessApiTestSuite) TestMerchantBusinessApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createMerchantBusiness(input: {
			merchant_id: %d
			business_type: "Retail"
			tax_id: "123-456-789"
			established_year: 2020
			number_of_employees: 10
			website_url: "http://example.com"
		}) {
			status message
			data { id business_type }
		}
	}`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createMerchantBusiness")
	s.Require().NotNil(data)
	s.Equal("Retail", data["business_type"])
	s.businessID = int(data["id"].(float64))
	s.Require().NotZero(s.businessID)

	// 2. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findMerchantBusinessById(input: { id: %d }) { status message data { id business_type } } }`, s.businessID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findMerchantBusinessById")
	s.Require().NotNil(data)
	s.Equal(float64(s.businessID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllMerchantBusinesses(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllMerchantBusinesses"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveMerchantBusinesses(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveMerchantBusinesses"))

	// 5. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateMerchantBusiness(input: {
			merchant_business_info_id: %d
			business_type: "Wholesale"
			tax_id: "987-654-321"
			established_year: 2021
			number_of_employees: 20
			website_url: "http://updated.com"
		}) {
			status message
			data { id business_type }
		}
	}`, s.businessID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateMerchantBusiness")
	s.Require().NotNil(data)
	s.Equal("Wholesale", data["business_type"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchantBusiness(input: { id: %d }) { status message data { id } } }`, s.businessID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedMerchantBusinesses(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedMerchantBusinesses"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreMerchantBusiness(input: { id: %d }) { status message data { id } } }`, s.businessID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchantBusiness(input: { id: %d }) { status message } }`, s.businessID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteMerchantBusinessPermanent(input: { id: %d }) { status message } }`, s.businessID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllMerchantBusinesses { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllMerchantBusinessesPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestMerchantBusinessApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessApiTestSuite))
}
