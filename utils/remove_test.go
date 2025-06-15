package utils_test

import (
	"strings"

	"github.com/0x726f6f6b6965/openapi-fmt/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func (suite *UtilsTestSuite) TestRemoveExtensions() {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(file)
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
	if strings.Contains(str, "x-") {
		suite.T().Errorf("Extensions were not removed correctly, found: %s", str)
	}
}
