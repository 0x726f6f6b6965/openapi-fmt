package utils_test

import (
	"github.com/0x726f6f6b6965/openapi-fmt/utils"
)

func (suite *UtilsTestSuite) TestSplitByPath() {
	output, err := utils.SplitByPath(suite.Doc, map[string]struct{}{
		"/api/v1/example": {}})
	if err != nil {
		suite.T().Fatal(err)
	}
	if _, exist := output.Components.Schemas["User"]; exist {
		suite.T().Errorf("Schema 'User' should not exist in the split document")
	}
	shouldExist := []string{
		"#/components/parameters/size",
		"#/components/schemas/JsonResponse",
		"#/components/responses/JsonResponse",
		"#/components/responses/HTTPResponseServiceUnavailable",
	}
	for _, ref := range shouldExist {
		typ, key := utils.ExtractReferenceName(ref)
		switch typ {
		case "schemas":
			if _, exist := output.Components.Schemas[key]; !exist {
				suite.T().Errorf("Schema '%s' should exist in the split document", key)
			}
		case "responses":
			if _, exist := output.Components.Responses[key]; !exist {
				suite.T().Errorf("Response '%s' should exist in the split document", key)
			}
		case "parameters":
			if _, exist := output.Components.Parameters[key]; !exist {
				suite.T().Errorf("Parameter '%s' should exist in the split document", key)
			}
		default:
			suite.T().Errorf("Unknown reference type '%s' in '%s'", typ, ref)
		}
	}
}
