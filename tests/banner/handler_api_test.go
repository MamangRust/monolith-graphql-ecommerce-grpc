package banner_test

import (
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type BannerApiTestSuite struct {
	tests.BaseTestSuite
	gql      http.Handler
	bannerID int
}

func (s *BannerApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupBannerService()

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

// do runs a GraphQL operation against the gateway and returns the decoded body.
func (s *BannerApiTestSuite) do(query string) map[string]any {
	body, err := tests.DoGraphQL(s.gql, query)
	s.Require().NoError(err)
	return body
}

func (s *BannerApiTestSuite) TestBannerApiLifecycle() {
	// 1. Create
	createQuery := `mutation {
		createBanner(input: {
			name: "Test Banner"
			start_date: "2024-01-01"
			end_date: "2024-12-31"
			start_time: "00:00:00"
			end_time: "23:59:59"
			is_active: true
		}) {
			status
			message
			data { banner_id name }
		}
	}`
	body := s.do(createQuery)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createBanner")
	s.Require().NotNil(data)
	s.Equal("Test Banner", data["name"])
	s.bannerID = int(data["banner_id"].(float64))
	s.Require().NotZero(s.bannerID)

	// 2. FindById
	findQuery := fmt.Sprintf(`query {
		findBannerById(input: { id: %d }) {
			status message
			data { banner_id name }
		}
	}`, s.bannerID)
	body = s.do(findQuery)
	s.Require().Empty(tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findBannerById")
	s.Require().NotNil(data)
	s.Equal(float64(s.bannerID), data["banner_id"])

	// 3. FindAll
	body = s.do(`query { findAllBanners(input: { page: 1, page_size: 10 }) { status message data { banner_id } pagination { total_records } } }`)
	s.Require().Empty(tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllBanners"))

	// 4. FindByActive
	body = s.do(`query { findActiveBanners(input: { page: 1, page_size: 10 }) { status message data { banner_id } } }`)
	s.Require().Empty(tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveBanners"))

	// 5. Update
	updateQuery := fmt.Sprintf(`mutation {
		updateBanner(input: {
			banner_id: %d
			name: "Updated Banner"
			start_date: "2024-01-01"
			end_date: "2024-12-31"
			start_time: "00:00:00"
			end_time: "23:59:59"
			is_active: false
		}) {
			status message
			data { banner_id name is_active }
		}
	}`, s.bannerID)
	body = s.do(updateQuery)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateBanner")
	s.Require().NotNil(data)
	s.Equal("Updated Banner", data["name"])
	s.Equal(false, data["is_active"])

	// 6. Trash
	trashQuery := fmt.Sprintf(`mutation { trashBanner(input: { id: %d }) { status message data { banner_id } } }`, s.bannerID)
	body = s.do(trashQuery)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body = s.do(`query { findTrashedBanners(input: { page: 1, page_size: 10 }) { status message data { banner_id } } }`)
	s.Require().Empty(tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedBanners"))

	// 8. Restore
	restoreQuery := fmt.Sprintf(`mutation { restoreBanner(input: { id: %d }) { status message data { banner_id } } }`, s.bannerID)
	body = s.do(restoreQuery)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent
	permQuery := fmt.Sprintf(`mutation { deleteBannerPermanent(input: { id: %d }) { status message } }`, s.bannerID)
	body = s.do(permQuery)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body = s.do(`mutation { restoreAllBanners { status message } }`)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body = s.do(`mutation { deleteAllBannersPermanent { status message } }`)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestBannerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerApiTestSuite))
}
