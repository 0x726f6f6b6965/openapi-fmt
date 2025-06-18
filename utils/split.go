package utils

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func SplitByPath(doc *openapi3.T, targets map[string]struct{}) (*openapi3.T, error) {
	if doc == nil {
		return nil, ErrOpenAPINotFound
	}

	// Create a new OpenAPI document to hold the split paths
	splitDoc := &openapi3.T{
		OpenAPI: doc.OpenAPI,
		Info:    doc.Info,
		Servers: doc.Servers,
	}

	// Initialize the Paths field
	splitDoc.Paths = openapi3.NewPaths()
	// Initialize the Components field
	splitDoc.Components = &openapi3.Components{
		Schemas:         make(openapi3.Schemas),
		Responses:       make(openapi3.ResponseBodies),
		Parameters:      make(openapi3.ParametersMap),
		RequestBodies:   make(openapi3.RequestBodies),
		Headers:         make(openapi3.Headers),
		SecuritySchemes: make(openapi3.SecuritySchemes),
	}

	for path, pathItem := range doc.Paths.Map() {
		if _, ok := targets[path]; !ok {
			continue
		}
		splitDoc.Paths.Set(path, pathItem)
		// Collect components that are referenced in the path item
		for _, operation := range pathItem.Operations() {
			if operation == nil {
				continue
			}
			// Collect parameters
			for _, param := range operation.Parameters {
				collectPrameterComponents(splitDoc, doc, param, "")
			}
			// Collect request body
			if operation.RequestBody != nil {
				collectRequestBodyComponents(splitDoc, doc, operation.RequestBody, "")
			}

			// Collect responses
			for _, responseRef := range operation.Responses.Map() {
				collectResponseComponents(splitDoc, doc, responseRef, "")
			}

		}

	}
	if len(splitDoc.Paths.Map()) == 0 {
		return nil, ErrOpenAPIPathNotFound
	}
	return splitDoc, nil
}

func collectPrameterComponents(splitDoc *openapi3.T, doc *openapi3.T, param *openapi3.ParameterRef, ref string) {
	if param == nil {
		return
	}
	if param.Ref != "" {
		// extract the reference name
		typ, key := ExtractReferenceName(param.Ref)
		switch typ {
		case "parameters":
			// If the reference is to a parameter, we need to collect it
			// from the main document's components
			next := doc.Components.Parameters[key]
			collectPrameterComponents(splitDoc, doc, next, param.Ref)
		case "schemas":
			// If the reference is to a schema, we need to collect it
			// from the main document's components
			next := doc.Components.Schemas[key]
			collectSchemaComponents(splitDoc, doc, next, param.Ref)
		default:
		}
		return
	}
	if ref != "" {
		// If the parameter is not a reference, we can add it directly
		typ, key := ExtractReferenceName(ref)
		if typ != "parameters" {
			// If the reference is not to a parameter, we should not add it
			return
		}
		if splitDoc.Components.Parameters == nil {
			splitDoc.Components.Parameters = make(openapi3.ParametersMap)
		}
		splitDoc.Components.Parameters[key] = param
	}
}

func collectRequestBodyComponents(splitDoc *openapi3.T, doc *openapi3.T, requestBody *openapi3.RequestBodyRef, ref string) {
	if requestBody == nil {
		return
	}
	if requestBody.Ref != "" {
		// extract the reference name
		typ, key := ExtractReferenceName(requestBody.Ref)
		switch typ {
		case "requestBodies":
			// If the reference is to a request body, we need to collect it
			// from the main document's components
			next := doc.Components.RequestBodies[key]
			collectRequestBodyComponents(splitDoc, doc, next, requestBody.Ref)
		case "schemas":
			// If the reference is to a schema, we need to collect it
			// from the main document's components
			next := doc.Components.Schemas[key]
			collectSchemaComponents(splitDoc, doc, next, requestBody.Ref)
		default:
		}
		return
	}

	if ref != "" {
		// If the request body is not a reference, we can add it directly
		// extract the reference name
		typ, key := ExtractReferenceName(ref)
		if typ != "requestBodies" {
			// If the reference is not to a request body, we should not add it
			return
		}
		// Add the request body to the split document's components
		if splitDoc.Components.RequestBodies == nil {
			splitDoc.Components.RequestBodies = make(openapi3.RequestBodies)
		}
		splitDoc.Components.RequestBodies[key] = requestBody
		return
	}
	if requestBody.Value != nil {
		for _, mediaItem := range requestBody.Value.Content {
			if mediaItem.Schema != nil {
				// Collect schema components from the media type schema
				collectSchemaComponents(splitDoc, doc, mediaItem.Schema, "")
			}
		}
	}
}

func collectResponseComponents(splitDoc *openapi3.T, doc *openapi3.T, response *openapi3.ResponseRef, ref string) {
	if response == nil {
		return
	}
	if response.Ref != "" {
		typ, key := ExtractReferenceName(response.Ref)
		switch typ {
		case "responses":
			next := doc.Components.Responses[key]
			collectResponseComponents(splitDoc, doc, next, response.Ref)
		case "schemas":
			next := doc.Components.Schemas[key]
			collectSchemaComponents(splitDoc, doc, next, response.Ref)
		default:
		}
		return
	}

	if ref != "" {
		// If the response is not a reference, we can add it directly
		typ, key := ExtractReferenceName(ref)
		if typ != "responses" {
			return
		}
		if splitDoc.Components.Responses == nil {
			splitDoc.Components.Responses = make(openapi3.ResponseBodies)
		}
		// Add the response to the split document's components
		splitDoc.Components.Responses[key] = response
		return
	}
	if response.Value != nil {
		for _, mediaItem := range response.Value.Content {
			if mediaItem.Schema != nil {
				// Collect schema components from the media type schema
				collectSchemaComponents(splitDoc, doc, mediaItem.Schema, "")
			}
		}
	}

}

func collectSchemaComponents(splitDoc *openapi3.T, doc *openapi3.T, schema *openapi3.SchemaRef, ref string) {
	if schema == nil {
		return
	}
	if schema.Ref != "" {
		_, key := ExtractReferenceName(schema.Ref)
		next := doc.Components.Schemas[key]
		collectSchemaComponents(splitDoc, doc, next, schema.Ref)
		return
	}
	if ref != "" {
		typ, key := ExtractReferenceName(ref)
		if typ != "schemas" {
			return
		}
		if splitDoc.Components.Schemas == nil {
			splitDoc.Components.Schemas = make(openapi3.Schemas)
		}
		splitDoc.Components.Schemas[key] = schema
	}
	// Collect schema components from the schema's properties
	for _, property := range schema.Value.Properties {
		collectSchemaComponents(splitDoc, doc, property, "")
	}
	// Collect schema components from the schema's items
	if schema.Value.Items != nil {
		collectSchemaComponents(splitDoc, doc, schema.Value.Items, "")
	}
	// Collect schema components from the schema's allOf, oneOf, anyOf
	for _, ref := range schema.Value.AllOf {
		collectSchemaComponents(splitDoc, doc, ref, "")
	}
	for _, ref := range schema.Value.OneOf {
		collectSchemaComponents(splitDoc, doc, ref, "")
	}
	for _, ref := range schema.Value.AnyOf {
		collectSchemaComponents(splitDoc, doc, ref, "")
	}
}

func ExtractReferenceName(ref string) (string, string) {
	// Split the reference string by '/' and take the last part
	parts := strings.Split(ref, "/")
	if len(parts) == 0 {
		return "", ""
	}
	return parts[len(parts)-2], parts[len(parts)-1]
}
