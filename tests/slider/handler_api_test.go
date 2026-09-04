package slider_test

import (
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type SliderApiTestSuite struct {
	tests.BaseTestSuite
	gql      http.Handler
	sliderID int
}

const sliderImageContent = "dummy image content"

func (s *SliderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupSliderService()

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *SliderApiTestSuite) TestSliderApiLifecycle() {
	// 1. Create (multipart upload for the required image field)
	body, err := tests.DoGraphQLWithUploads(
		s.gql,
		`mutation($input: CreateSliderRequest!) {
			createSlider(input: $input) { status message data { id name image } }
		}`,
		map[string]any{"input": map[string]any{"name": "Test Slider", "image": nil}},
		[]tests.GraphQLUpload{{VarPath: "variables.input.image", Filename: "slider.jpg", Content: []byte(sliderImageContent)}},
	)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createSlider")
	s.Require().NotNil(data)
	s.Equal("Test Slider", data["name"])
	s.sliderID = int(data["id"].(float64))
	s.Require().NotZero(s.sliderID)

	// 2. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllSliders(input: { page: 1, page_size: 10 }) { status message data { id name } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllSliders"))

	// 3. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveSliders(input: { page: 1, page_size: 10 }) { status message data { id name } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveSliders"))

	// 4. Update
	body, err = tests.DoGraphQLWithUploads(
		s.gql,
		`mutation($input: UpdateSliderRequest!) {
			updateSlider(input: $input) { status message data { id name } }
		}`,
		map[string]any{"input": map[string]any{
			"id":    float64(s.sliderID),
			"name":  "Updated Test Slider",
			"image": nil,
		}},
		[]tests.GraphQLUpload{{VarPath: "variables.input.image", Filename: "slider_updated.jpg", Content: []byte(sliderImageContent)}},
	)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateSlider")
	s.Require().NotNil(data)
	s.Equal("Updated Test Slider", data["name"])

	// 5. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedSlider(input: { id: %d }) { status message data { id } } }`, s.sliderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 6. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedSliders(input: { page: 1, page_size: 10 }) { status message data { id name } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedSliders"))

	// 7. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreSlider(input: { id: %d }) { status message data { id } } }`, s.sliderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 8. DeletePermanent (trash it again first, like the old REST flow did)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedSlider(input: { id: %d }) { status message } }`, s.sliderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteSliderPermanent(input: { id: %d }) { status message } }`, s.sliderID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllSliders { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllSlidersPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestSliderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderApiTestSuite))
}
