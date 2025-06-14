package utils_test

import (
	"strings"
	"testing"

	_ "embed"

	"github.com/0x726f6f6b6965/openapi-fmt/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/api.yaml
var file []byte

func TestRemoveExtensions(t *testing.T) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(file)
	if err != nil {
		t.Fatal(err)
	}

	exclude := map[string]struct{}{}

	utils.RemoveExtensions(doc, exclude)
	f, err := doc.MarshalYAML()
	if err != nil {
		t.Fatal(err)
	}
	b, err := yaml.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	// Check if the extensions were removed correctly
	str := string(b)
	if strings.Contains(str, "x-") {
		t.Errorf("Extensions were not removed correctly, found: %s", str)
	}
}
