package utils_test

import (
	"github.com/0x726f6f6b6965/openapi-fmt/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func (suite *UtilsTestSuite) TestRemoveExtensions() {
	loader := openapi3.NewLoader()
	// Assuming 'file' is defined elsewhere in your actual test setup
	// For this example, let's use a simple OpenAPI spec string
	fileContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test endpoint
      x-go-type: "TestType"
      x-42c-sample: "SampleValue"
      responses:
        '200':
          description: Successful response
`
	doc, err := loader.LoadFromData([]byte(fileContent))
	if err != nil {
		suite.T().Fatal(err)
	}

	exclude := map[string]struct{}{}

	utils.RemoveExtensions(doc, exclude)
	f, err := doc.MarshalYAML()
	if err != nil {
		suite.T().Fatal(err)
	}
	b, err := yaml.Marshal(f)
	if err != nil {
		suite.T().Fatal(err)
	}
	// Check if the extensions were removed correctly
	str := string(b)
	assert.NotContains(suite.T(), str, "x-", "Extensions were not removed correctly")
}

func (suite *UtilsTestSuite) TestRemoveExtensionsWithExcludes() {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(removeExtensionsFile)
	if err != nil {
		suite.T().Fatal(err)
	}

	exclude := map[string]struct{}{
		"x-go-type": {},
	}

	utils.RemoveExtensions(doc, exclude)
	f, err := doc.MarshalYAML()
	if err != nil {
		suite.T().Fatal(err)
	}
	b, err := yaml.Marshal(f)
	if err != nil {
		suite.T().Fatal(err)
	}

	str := string(b)
	assert.Contains(suite.T(), str, "x-go-type", "Excluded extension 'x-go-type' should be present")
	assert.NotContains(suite.T(), str, "x-42c-sample", "Non-excluded extension 'x-42c-sample' should not be present")
}

func (suite *UtilsTestSuite) TestRemoveExtensionsNoExtensionsPresent() {
	loader := openapi3.NewLoader()
	minimalContent := `
openapi: 3.0.0
info:
  title: Minimal API
  version: 0.1.0
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: OK
`
	doc, err := loader.LoadFromData([]byte(minimalContent))
	if err != nil {
		suite.T().Fatal(err)
	}

	originalYAML, err := yaml.Marshal(doc)
	if err != nil {
		suite.T().Fatal(err)
	}

	utils.RemoveExtensions(doc, map[string]struct{}{})
	processedYAMLNode, err := doc.MarshalYAML()
	if err != nil {
		suite.T().Fatal(err)
	}
	processedYAML, err := yaml.Marshal(processedYAMLNode)
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.YAMLEq(suite.T(), string(originalYAML), string(processedYAML), "YAML output should be unchanged")
	assert.NotContains(suite.T(), string(processedYAML), "x-", "Output should not contain any x- extensions")
}

func (suite *UtilsTestSuite) TestRemoveExtensionsNilDocument() {
	assert.NotPanics(suite.T(), func() {
		utils.RemoveExtensions(nil, map[string]struct{}{})
	}, "RemoveExtensions should not panic with nil document and empty exclude map")

	assert.NotPanics(suite.T(), func() {
		utils.RemoveExtensions(nil, map[string]struct{}{"x-keep": {}})
	}, "RemoveExtensions should not panic with nil document and non-empty exclude map")
}
