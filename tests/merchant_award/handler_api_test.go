package merchant_award_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantAwardApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	awardID    int
	merchantID int
}

func (s *MerchantAwardApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupMerchantAwardService()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, userID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *MerchantAwardApiTestSuite) TestMerchantAwardApiLifecycle() {
	// 1. Create
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createMerchantAward(input: {
			merchant_id: %d
			title: "Best Merchant 2024"
			description: "Award for excellence"
			issued_by: "E-commerce Platform"
			issue_date: "2024-01-01"
			expiry_date: "2025-01-01"
			certificate_url: "http://example.com/cert.pdf"
		}) {
			status message
			data { id title }
		}
	}`, s.merchantID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createMerchantAward")
	s.Require().NotNil(data)
	s.Equal("Best Merchant 2024", data["title"])
	s.awardID = int(data["id"].(float64))
	s.Require().NotZero(s.awardID)

	// 2. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findMerchantAwardById(input: { id: %d }) { status message data { id title } } }`, s.awardID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findMerchantAwardById")
	s.Require().NotNil(data)
	s.Equal(float64(s.awardID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllMerchantAwards(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllMerchantAwards"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveMerchantAwards(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveMerchantAwards"))

	// 5. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateMerchantAward(input: {
			merchant_certification_id: %d
			title: "Updated Award Title"
			description: "Updated Description"
			issued_by: "Updated Issuer"
			issue_date: "2024-02-01"
			expiry_date: "2025-02-01"
			certificate_url: "http://example.com/updated.pdf"
		}) {
			status message
			data { id title }
		}
	}`, s.awardID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateMerchantAward")
	s.Require().NotNil(data)
	s.Equal("Updated Award Title", data["title"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchantAward(input: { id: %d }) { status message data { id } } }`, s.awardID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedMerchantAwards(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedMerchantAwards"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreMerchantAward(input: { id: %d }) { status message data { id } } }`, s.awardID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchantAward(input: { id: %d }) { status message } }`, s.awardID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteMerchantAwardPermanent(input: { id: %d }) { status message } }`, s.awardID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllMerchantAwards { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllMerchantAwardsPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestMerchantAwardApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardApiTestSuite))
}
