package category_test

import (
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-graphql-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CategoryApiTestSuite struct {
	tests.BaseTestSuite
	gql        http.Handler
	categoryID int
}

const categoryImageContent = "dummy image content"

func (s *CategoryApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupCategoryService()

	s.gql = tests.NewGraphQLHandler(s.Conns, s.GetCacheStore(), s.Log)
}

func (s *CategoryApiTestSuite) createCategory(name string) (float64, map[string]any) {
	body, err := tests.DoGraphQLWithUploads(
		s.gql,
		`mutation($input: CreateCategoryInput!) {
			createCategory(input: $input) { status message data { id name description slug_category image_category } }
		}`,
		map[string]any{"input": map[string]any{
			"name":           name,
			"description":    "Test Description",
			"slug_category":  "test-category",
			"image_category": nil,
		}},
		[]tests.GraphQLUpload{{VarPath: "variables.input.image_category", Filename: "test.jpg", Content: []byte(categoryImageContent)}},
	)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createCategory")
	s.Require().NotNil(data)
	s.Equal(name, data["name"])
	id, ok := data["id"].(float64)
	s.Require().True(ok)
	s.Require().NotZero(int(id))
	return id, data
}

func (s *CategoryApiTestSuite) TestCategoryApiLifecycle() {
	// 1. Create
	categoryID, _ := s.createCategory("Test Category")
	s.categoryID = int(categoryID)

	// 2. FindById
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`query {
		findCategoryById(input: { id: %d }) { status message data { id name } }
	}`, s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "findCategoryById")
	s.Require().NotNil(data)
	s.Equal(float64(s.categoryID), data["id"])

	// 3. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllCategories(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllCategories"))

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findActiveCategories(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findActiveCategories"))

	// 5. Update
	body, err = tests.DoGraphQLWithUploads(
		s.gql,
		`mutation($input: UpdateCategoryInput!) {
			updateCategory(input: $input) { status message data { id name description } }
		}`,
		map[string]any{"input": map[string]any{
			"category_id":    float64(s.categoryID),
			"name":           "Updated Category",
			"description":    "Updated Description",
			"slug_category":  "updated-category",
			"image_category": nil,
		}},
		[]tests.GraphQLUpload{{VarPath: "variables.input.image_category", Filename: "updated.jpg", Content: []byte(categoryImageContent)}},
	)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateCategory")
	s.Require().NotNil(data)
	s.Equal("Updated Category", data["name"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashCategory(input: { id: %d }) { status message data { id } } }`, s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findTrashedCategories(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findTrashedCategories"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreCategory(input: { id: %d }) { status message data { id } } }`, s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashCategory(input: { id: %d }) { status message } }`, s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteCategoryPermanent(input: { id: %d }) { status message } }`, s.categoryID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllCategories { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllCategoriesPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestCategoryApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryApiTestSuite))
}
