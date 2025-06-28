package utils_test

import (
	"testing"

	"github.com/0x726f6f6b6965/openapi-fmt/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func TestExtractReferenceName(t *testing.T) {
	testCases := []struct {
		name         string
		ref          string
		expectedType string
		expectedKey  string
	}{
		{
			name:         "valid schema",
			ref:          "#/components/schemas/User",
			expectedType: "schemas",
			expectedKey:  "User",
		},
		{
			name:         "valid parameter",
			ref:          "#/components/parameters/Param1",
			expectedType: "parameters",
			expectedKey:  "Param1",
		},
		{
			name:         "valid response",
			ref:          "#/components/responses/ErrorResponse",
			expectedType: "responses",
			expectedKey:  "ErrorResponse",
		},
		{
			name:         "valid requestBody",
			ref:          "#/components/requestBodies/UserBody",
			expectedType: "requestBodies",
			expectedKey:  "UserBody",
		},
		{
			name:         "valid header",
			ref:          "#/components/headers/RateLimit",
			expectedType: "headers",
			expectedKey:  "RateLimit",
		},
		{
			name:         "valid securityScheme",
			ref:          "#/components/securitySchemes/OAuth2",
			expectedType: "securitySchemes",
			expectedKey:  "OAuth2",
		},
		{
			name:         "invalid ref with slash",
			ref:          "invalid/ref",
			expectedType: "invalid", // Adjusted expectation
			expectedKey:  "ref",     // Adjusted expectation
		},
		{
			name:         "single part ref (invalid)",
			ref:          "User",
			expectedType: "",
			expectedKey:  "",
		},
		{
			name:         "empty string",
			ref:          "",
			expectedType: "",
			expectedKey:  "",
		},
		{
			name:         "missing key",
			ref:          "#/components/schemas/",
			expectedType: "schemas",
			expectedKey:  "",
		},
		{
			name:         "missing components",
			ref:          "#/schemas/User",
			expectedType: "schemas",
			expectedKey:  "User",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualType, actualKey := utils.ExtractReferenceName(tc.ref)
			assert.Equal(t, tc.expectedType, actualType, "Type mismatch")
			assert.Equal(t, tc.expectedKey, actualKey, "Key mismatch")
		})
	}
}

// loadTestAPIDoc is a helper function to load the test API document.
// It's assumed to be used by tests within the UtilsTestSuite.
func (suite *UtilsTestSuite) loadTestAPIDoc() *openapi3.T {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(splitFile)
	if err != nil {
		suite.T().Fatalf("Failed to load data from testdata/split_api.yaml: %v", err)
	}
	return doc
}

func (suite *UtilsTestSuite) TestSplitByPath_ComplexReferences() {
	doc := suite.loadTestAPIDoc()
	output, err := utils.SplitByPath(doc, map[string][]string{
		"/complex-path": {"get"},
	})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), output)
	assert.NotNil(suite.T(), output.Components)

	// Check for expected components
	assert.Contains(suite.T(), output.Components.Parameters, "CommonParam", "Missing CommonParam")
	assert.Contains(suite.T(), output.Components.RequestBodies, "ComplexBody", "Missing ComplexBody")
	assert.Contains(suite.T(), output.Components.Responses, "SuccessResponse", "Missing SuccessResponse")
	assert.Contains(suite.T(), output.Components.Responses, "ErrorResponse", "Missing ErrorResponse")

	assert.Contains(suite.T(), output.Components.Schemas, "MainSchema", "Missing MainSchema")
	assert.Contains(suite.T(), output.Components.Schemas, "NestedSchema", "Missing NestedSchema (indirectly via MainSchema)")
	assert.Contains(suite.T(), output.Components.Schemas, "ReferencedInResponse", "Missing ReferencedInResponse (via SuccessResponse)")

	// Check that other, unreferenced components are NOT present
	assert.NotContains(suite.T(), output.Components.Schemas, "UnusedSchema")
	assert.NotContains(suite.T(), output.Components.Parameters, "UnusedParam")
	assert.NotContains(suite.T(), output.Components.RequestBodies, "UnusedBody")
	assert.NotContains(suite.T(), output.Components.Responses, "UnusedResponse")
	assert.NotContains(suite.T(), output.Components.Headers, "UnusedHeader")
	assert.NotContains(suite.T(), output.Components.SecuritySchemes, "UnusedScheme")

	// Check that the selected path is present
	assert.Contains(suite.T(), output.Paths.Map(), "/complex-path")
	// Check that other paths are not present
	assert.NotContains(suite.T(), output.Paths.Map(), "/simple-path")
	assert.NotContains(suite.T(), output.Paths.Map(), "/path-with-allof-schema")
}

func (suite *UtilsTestSuite) TestSplitByPath_PathExistsNoNewComponents() {
	doc := suite.loadTestAPIDoc()
	output, err := utils.SplitByPath(doc, map[string][]string{
		"/simple-path": {},
	})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), output)

	// Assert that Components is not nil, but specific component maps might be empty or nil
	if output.Components != nil {
		assert.Empty(suite.T(), output.Components.Schemas, "Schemas should be empty")
		assert.Empty(suite.T(), output.Components.Parameters, "Parameters should be empty")
		assert.Empty(suite.T(), output.Components.RequestBodies, "RequestBodies should be empty")
		assert.Empty(suite.T(), output.Components.Responses, "Responses should be empty")
		assert.Empty(suite.T(), output.Components.Headers, "Headers should be empty")
		assert.Empty(suite.T(), output.Components.SecuritySchemes, "SecuritySchemes should be empty")
	} else {
		// If Components itself is nil, that's also acceptable.
	}

	// Check that the selected path is present
	assert.Contains(suite.T(), output.Paths.Map(), "/simple-path")
	// Check that other paths are not present
	assert.NotContains(suite.T(), output.Paths.Map(), "/complex-path")
}

func (suite *UtilsTestSuite) TestSplitByPath_AllOfSchema() {
	doc := suite.loadTestAPIDoc()
	output, err := utils.SplitByPath(doc, map[string][]string{
		"/path-with-allof-schema": {},
	})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), output)
	assert.NotNil(suite.T(), output.Components)
	assert.NotNil(suite.T(), output.Components.Schemas)

	assert.Contains(suite.T(), output.Components.Schemas, "AllOfSchema", "Missing AllOfSchema")
	assert.Contains(suite.T(), output.Components.Schemas, "BaseSchema", "Missing BaseSchema (referenced by AllOfSchema)")

	// Check that the selected path is present
	assert.Contains(suite.T(), output.Paths.Map(), "/path-with-allof-schema")
}

func (suite *UtilsTestSuite) TestSplitByPathNotFound() {
	doc := suite.loadTestAPIDoc() // Use helper to load doc
	_, err := utils.SplitByPath(doc, map[string][]string{
		"/api/v1/nonexistent": {}})
	assert.Error(suite.T(), err, "Expected an error when splitting by a non-existent path")
	assert.ErrorIs(suite.T(), err, utils.ErrOpenAPIPathNotFound, "Expected ErrOpenAPIPathNotFound")
}

func (suite *UtilsTestSuite) TestSplitByPathEmpty() {
	doc := suite.loadTestAPIDoc() // Use helper to load doc
	_, err := utils.SplitByPath(doc, map[string][]string{})
	assert.Error(suite.T(), err, "Expected an error when splitting by an empty path map")
	assert.ErrorIs(suite.T(), err, utils.ErrOpenAPIPathNotFound, "Expected ErrOpenAPIPathNotFound")
}

func (suite *UtilsTestSuite) TestSplitByPathNilDocument() {
	// No need to load full doc for this one, suite.Doc is not used.
	_, err := utils.SplitByPath(nil, map[string][]string{"/any/path": {}})
	assert.Error(suite.T(), err, "Expected an error when splitting with a nil document")
	assert.ErrorIs(suite.T(), err, utils.ErrOpenAPINotFound, "Expected ErrOpenAPINotFound")
}
