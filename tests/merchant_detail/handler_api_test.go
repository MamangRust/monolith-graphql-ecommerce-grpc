package merchant_detail_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantDetailApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	detailID   int
	merchantID int
	userID     int
}

const merchantDetailImageContent = "dummy image content"

func (s *MerchantDetailApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupMerchantDetailService()

	// Seed dependencies
	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *MerchantDetailApiTestSuite) TestMerchantDetailApiLifecycle() {
	// 1. Create (cover_image and logo are Upload fields)
	body, err := tests.DoGraphQLWithUploads(
		s.gql,
		`mutation($input: CreateMerchantDetailInput!) {
			createMerchantDetail(input: $input) { status message data { id display_name } }
		}`,
		map[string]any{"input": map[string]any{
			"merchant_id":       float64(s.merchantID),
			"display_name":      "Test Merchant",
			"cover_image":       nil,
			"logo":              nil,
			"short_description": "A test merchant",
			"website_url":       "http://example.com",
		}},
		[]tests.GraphQLUpload{
			{VarPath: "variables.input.cover_image", Filename: "cover.jpg", Content: []byte(merchantDetailImageContent)},
			{VarPath: "variables.input.logo", Filename: "logo.jpg", Content: []byte(merchantDetailImageContent)},
		},
	)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createMerchantDetail")
	s.Require().NotNil(data)
	s.Equal("Test Merchant", data["display_name"])
	s.detailID = int(data["id"].(float64))
	s.Require().NotZero(s.detailID)

	// 2. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findMerchantDetailById(input: { id: %d }) { status message data { id display_name } } }`, s.detailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findMerchantDetailById")
	s.Require().NotNil(data)
	s.Equal(float64(s.detailID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllMerchantDetails(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllMerchantDetails"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveMerchantDetails(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveMerchantDetails"))

	// 5. Update
	body, err = tests.DoGraphQLWithUploads(
		s.gql,
		`mutation($input: UpdateMerchantDetailInput!) {
			updateMerchantDetail(input: $input) { status message data { id display_name } }
		}`,
		map[string]any{"input": map[string]any{
			"merchant_detail_id": float64(s.detailID),
			"display_name":       "Updated Merchant",
			"cover_image":        nil,
			"logo":               nil,
			"short_description":  "Updated short description",
			"website_url":        "http://updated.com",
		}},
		[]tests.GraphQLUpload{
			{VarPath: "variables.input.cover_image", Filename: "cover_updated.jpg", Content: []byte(merchantDetailImageContent)},
			{VarPath: "variables.input.logo", Filename: "logo_updated.jpg", Content: []byte(merchantDetailImageContent)},
		},
	)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateMerchantDetail")
	s.Require().NotNil(data)
	s.Equal("Updated Merchant", data["display_name"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchantDetail(input: { id: %d }) { status message data { id } } }`, s.detailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedMerchantDetails(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedMerchantDetails"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreMerchantDetail(input: { id: %d }) { status message data { id } } }`, s.detailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashMerchantDetail(input: { id: %d }) { status message } }`, s.detailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteMerchantDetailPermanent(input: { id: %d }) { status message } }`, s.detailID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllMerchantDetails { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllMerchantDetailsPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestMerchantDetailApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailApiTestSuite))
}
