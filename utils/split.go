package utils

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func SplitByPath(doc *openapi3.T, targets map[string][]string) (*openapi3.T, error) {
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
		op := make(map[string]bool)
		for _, method := range targets[path] {
			if method == "" {
				continue
			}
			// Normalize method to uppercase
			op[strings.ToUpper(method)] = true
		}
		var allOperations bool
		if len(op) == 0 {
			allOperations = true // If no specific methods are provided, include all operations
		}
		splitDoc.Paths.Set(path, pathItem)
		// Collect components that are referenced in the path item
		for method, operation := range pathItem.Operations() {
			if operation == nil {
				continue
			}
			if !allOperations && len(op) > 0 && !op[method] {
				continue // Skip operations not in the specified methods
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
		typ, key := ExtractReferenceName(param.Ref)
		switch typ {
		case "parameters":
			if _, exists := splitDoc.Components.Parameters[key]; exists {
				return
			}
			sourceComponent := doc.Components.Parameters[key]
			if sourceComponent == nil {
				return
			}

			if splitDoc.Components.Parameters == nil {
				splitDoc.Components.Parameters = make(openapi3.ParametersMap)
			}
			splitDoc.Components.Parameters[key] = sourceComponent
			collectPrameterComponents(splitDoc, doc, sourceComponent, param.Ref) // Recurse with source component and original ref
		default:
		}
		return
	}

	// Logic for handling resolved parameter 'param' (param.Ref is empty)
	// 'ref' is the original reference string that led to this 'param' value.
	if ref != "" {
		originalRefType, originalRefKey := ExtractReferenceName(ref)
		if originalRefType == "parameters" {
			if _, exists := splitDoc.Components.Parameters[originalRefKey]; !exists {
				if splitDoc.Components.Parameters == nil {
					splitDoc.Components.Parameters = make(openapi3.ParametersMap)
				}
				splitDoc.Components.Parameters[originalRefKey] = param // Add the fully resolved parameter
			}
		}
	}
	// If param.Value is not nil and contains an inline schema, collect it.
	if param.Value != nil && param.Value.Schema != nil {
		collectSchemaComponents(splitDoc, doc, param.Value.Schema)
	}
}

func collectRequestBodyComponents(splitDoc *openapi3.T, doc *openapi3.T, requestBody *openapi3.RequestBodyRef, ref string) {
	if requestBody == nil {
		return
	}
	if requestBody.Ref != "" {
		typ, key := ExtractReferenceName(requestBody.Ref)
		switch typ {
		case "requestBodies":
			if _, exists := splitDoc.Components.RequestBodies[key]; exists {
				return
			}
			sourceComponent := doc.Components.RequestBodies[key]
			if sourceComponent == nil {
				return
			}

			if splitDoc.Components.RequestBodies == nil {
				splitDoc.Components.RequestBodies = make(openapi3.RequestBodies)
			}
			splitDoc.Components.RequestBodies[key] = sourceComponent
			collectRequestBodyComponents(splitDoc, doc, sourceComponent, requestBody.Ref)
		default:
		}
		return
	}

	if ref != "" {
		originalRefType, originalRefKey := ExtractReferenceName(ref)
		if originalRefType == "requestBodies" {
			if _, exists := splitDoc.Components.RequestBodies[originalRefKey]; !exists {
				if splitDoc.Components.RequestBodies == nil {
					splitDoc.Components.RequestBodies = make(openapi3.RequestBodies)
				}
				splitDoc.Components.RequestBodies[originalRefKey] = requestBody
			}
		}
	}

	if requestBody.Value != nil {
		for _, mediaItem := range requestBody.Value.Content {
			if mediaItem.Schema != nil {
				collectSchemaComponents(splitDoc, doc, mediaItem.Schema)
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
			if _, exists := splitDoc.Components.Responses[key]; exists {
				return
			}
			sourceComponent := doc.Components.Responses[key]
			if sourceComponent == nil {
				return
			}

			if splitDoc.Components.Responses == nil {
				splitDoc.Components.Responses = make(openapi3.ResponseBodies)
			}
			splitDoc.Components.Responses[key] = sourceComponent
			collectResponseComponents(splitDoc, doc, sourceComponent, response.Ref)
		default:
		}
		return
	}

	if ref != "" {
		originalRefType, originalRefKey := ExtractReferenceName(ref)
		if originalRefType == "responses" {
			if _, exists := splitDoc.Components.Responses[originalRefKey]; !exists {
				if splitDoc.Components.Responses == nil {
					splitDoc.Components.Responses = make(openapi3.ResponseBodies)
				}
				splitDoc.Components.Responses[originalRefKey] = response
			}
		}
	}

	if response.Value != nil {
		for _, mediaItem := range response.Value.Content {
			if mediaItem.Schema != nil {
				collectSchemaComponents(splitDoc, doc, mediaItem.Schema)
			}
		}
	}
}

func collectSchemaComponents(splitDoc *openapi3.T, doc *openapi3.T, schemaRefToProcess *openapi3.SchemaRef) {
	if schemaRefToProcess == nil {
		return
	}

	var currentSchemaValue *openapi3.Schema

	if schemaRefToProcess.Ref != "" { // It's a reference e.g. "#/components/schemas/MySchema"
		refType, refKey := ExtractReferenceName(schemaRefToProcess.Ref)
		if refType != "schemas" {
			return // Not a schema component reference
		}

		// If already added, get its Value to process children. Avoids re-adding/overwriting.
		if existingRef, exists := splitDoc.Components.Schemas[refKey]; exists {
			currentSchemaValue = existingRef.Value
		} else {
			sourceComponent := doc.Components.Schemas[refKey]
			if sourceComponent == nil {
				return // Source component not found
			}

			if splitDoc.Components.Schemas == nil {
				splitDoc.Components.Schemas = make(openapi3.Schemas)
			}
			splitDoc.Components.Schemas[refKey] = sourceComponent // Add the component (SchemaRef) from source doc
			currentSchemaValue = sourceComponent.Value            // Process the value of this newly added component
		}
	} else {
		// It's an inline schema or a pre-resolved schema value (e.g. from a previous step)
		currentSchemaValue = schemaRefToProcess.Value
	}

	if currentSchemaValue == nil {
		return // No actual schema content to process
	}

	// Process children of currentSchemaValue
	for _, propertySchemaRef := range currentSchemaValue.Properties {
		collectSchemaComponents(splitDoc, doc, propertySchemaRef)
	}
	if currentSchemaValue.Items != nil {
		collectSchemaComponents(splitDoc, doc, currentSchemaValue.Items)
	}
	for _, allOfSchemaRef := range currentSchemaValue.AllOf {
		collectSchemaComponents(splitDoc, doc, allOfSchemaRef)
	}
	for _, oneOfSchemaRef := range currentSchemaValue.OneOf {
		collectSchemaComponents(splitDoc, doc, oneOfSchemaRef)
	}
	for _, anyOfSchemaRef := range currentSchemaValue.AnyOf {
		collectSchemaComponents(splitDoc, doc, anyOfSchemaRef)
	}
}

func ExtractReferenceName(ref string) (string, string) {
	if ref == "" {
		return "", ""
	}
	// Split the reference string by '/'
	parts := strings.Split(ref, "/")
	// We need at least two parts (e.g., "schemas/User" or "#/components/schemas/User")
	// For "#/components/schemas/User", parts are ["#", "components", "schemas", "User"] (len 4)
	// For "schemas/User", parts are ["schemas", "User"] (len 2)
	// For "User" (invalid short ref), parts are ["User"] (len 1)
	if len(parts) < 2 {
		return "", "" // Or handle as an error, or return the single part as key if appropriate
	}
	// The component type is the second to last part, and the key is the last part.
	return parts[len(parts)-2], parts[len(parts)-1]
}
