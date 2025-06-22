package utils

import "github.com/getkin/kin-openapi/openapi3"

// RemoveExtensions removes all extensions from the OpenAPI document,
// except those specified in the exclude map.
// The exclude map contains keys of extensions that should not be removed.
func RemoveExtensions(doc *openapi3.T, exclude map[string]struct{}) {
	if doc == nil {
		return
	}

	// Remove extensions from the OpenAPI document itself
	removeExt(doc.Extensions, exclude)

	// Remove components extensions
	removeComponentsExtensions(doc.Components, exclude)

	// Remove paths extensions
	doc.Paths.Extensions = nil
	for _, pathItem := range doc.Paths.Map() {
		pathItem.Extensions = nil
		for _, operation := range pathItem.Operations() {
			// Remove extensions from each operation
			removeExt(operation.Extensions, exclude) // Corrected: use removeExt helper

			// Remove parameters from the operation
			for _, param := range operation.Parameters {
				removeParameterExtensions(param, exclude)
			}

			// Remove request body extensions
			if operation.RequestBody != nil {
				removeRequestBodyExtensions(operation.RequestBody, exclude)
			}

			// Remove responses extensions
			for _, response := range operation.Responses.Map() {
				removeResponseExtensions(response, exclude)
			}
		}
	}
}

func removeExt(ext map[string]any, exclude map[string]struct{}) {
	if ext == nil {
		return
	}
	for key := range ext {
		if _, ok := exclude[key]; !ok {
			delete(ext, key)
		}
	}
}

func removeComponentsExtensions(components *openapi3.Components, exclude map[string]struct{}) {
	if components == nil {
		return
	}

	for _, schema := range components.Schemas {
		removeSchemaExtensions(schema, exclude)
	}

	for _, response := range components.Responses {
		removeResponseExtensions(response, exclude)
	}

	for _, parameter := range components.Parameters {
		removeParameterExtensions(parameter, exclude)
	}

	for _, requestBody := range components.RequestBodies {
		removeRequestBodyExtensions(requestBody, exclude)
	}

	for _, header := range components.Headers {
		removeHeaderExtensions(header, exclude)
	}

	for _, securityScheme := range components.SecuritySchemes {
		removeSecuritySchemeExtensions(securityScheme, exclude)
	}
}
func removeSchemaExtensions(schema *openapi3.SchemaRef, exclude map[string]struct{}) {
	if schema == nil {
		return
	}
	// Remove extensions from the schema reference itself
	removeExt(schema.Extensions, exclude)

	if schema.Value == nil {
		return
	}
	// Remove extensions from the schema itself
	removeExt(schema.Value.Extensions, exclude)

	// Recursively remove extensions from properties
	for _, prop := range schema.Value.Properties {
		removeSchemaExtensions(prop, exclude)
	}

	// Remove extensions from items if it's an array
	if schema.Value.Items != nil {
		removeSchemaExtensions(schema.Value.Items, exclude)
	}

	// Remove extensions from allOf, oneOf, anyOf
	for _, ref := range schema.Value.AllOf {
		removeSchemaExtensions(ref, exclude)
	}
	for _, ref := range schema.Value.OneOf {
		removeSchemaExtensions(ref, exclude)
	}
	for _, ref := range schema.Value.AnyOf {
		removeSchemaExtensions(ref, exclude)
	}
}
func removeResponseExtensions(response *openapi3.ResponseRef, exclude map[string]struct{}) {
	if response == nil {
		return
	}
	// Remove extensions from the response reference itself
	removeExt(response.Extensions, exclude)

	if response.Value == nil {
		return
	}
	// Remove extensions from the response itself
	removeExt(response.Value.Extensions, exclude)

	// Remove extensions from headers
	for _, header := range response.Value.Headers {
		removeHeaderExtensions(header, exclude)
	}

	// Remove extensions from content
	for _, mediaType := range response.Value.Content {
		removeSchemaExtensions(mediaType.Schema, exclude)
	}
}
func removeHeaderExtensions(header *openapi3.HeaderRef, exclude map[string]struct{}) {
	if header == nil {
		return
	}
	// Remove extensions from the header reference itself
	removeExt(header.Extensions, exclude)

	if header.Value == nil {
		return
	}
	// Remove extensions from the header itself
	header.Value.Extensions = nil

	// Remove extensions from schema
	removeSchemaExtensions(header.Value.Schema, exclude)
}
func removeParameterExtensions(parameter *openapi3.ParameterRef, exclude map[string]struct{}) {
	if parameter == nil {
		return
	}
	// Remove extensions from the parameter reference itself
	removeExt(parameter.Extensions, exclude)

	if parameter.Value == nil {
		return
	}
	// Remove extensions from the parameter itself
	removeExt(parameter.Value.Extensions, exclude)

	// Remove extensions from schema
	// removeSchemaExtensions(parameter.Value.Schema, exclude)
}
func removeRequestBodyExtensions(requestBody *openapi3.RequestBodyRef, exclude map[string]struct{}) {
	if requestBody == nil {
		return
	}
	// Remove extensions from the request body reference itself
	removeExt(requestBody.Extensions, exclude)

	if requestBody.Value == nil {
		return
	}
	// Remove extensions from the request body itself
	removeExt(requestBody.Value.Extensions, exclude)

	// Remove extensions from content
	for _, mediaType := range requestBody.Value.Content {
		removeSchemaExtensions(mediaType.Schema, exclude)
	}
}
func removeSecuritySchemeExtensions(securityScheme *openapi3.SecuritySchemeRef, exclude map[string]struct{}) {
	if securityScheme == nil {
		return
	}
	// Remove extensions from the security scheme reference itself
	removeExt(securityScheme.Extensions, exclude)

	if securityScheme.Value == nil {
		return
	}
	// Remove extensions from the security scheme itself
	removeExt(securityScheme.Value.Extensions, exclude)
}
