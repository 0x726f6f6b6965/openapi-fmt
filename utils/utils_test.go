package utils_test

import (
	_ "embed"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/suite"
)

//go:embed testdata/api.yaml
var file []byte

type UtilsTestSuite struct {
	suite.Suite
	Doc *openapi3.T
}

func (suite *UtilsTestSuite) SetupTest() {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(file)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.Doc = doc
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
