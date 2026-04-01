package config

import (
	"fmt"
	"math"
	"slices"
)

var supportedScalarTypes = []string{"string", "integer", "number", "boolean"}

func ValidateInputSchema(schema map[string]any) error {
	if len(schema) == 0 {
		return nil
	}

	rootType, ok := schema["type"].(string)
	if !ok || rootType != "object" {
		return fmt.Errorf("input_schema.type must be %q", "object")
	}

	propertiesValue, ok := schema["properties"]
	if ok {
		properties, ok := propertiesValue.(map[string]any)
		if !ok {
			return fmt.Errorf("input_schema.properties must be an object")
		}

		for name, rawProperty := range properties {
			property, ok := rawProperty.(map[string]any)
			if !ok {
				return fmt.Errorf("input_schema.properties.%s must be an object", name)
			}

			propertyType, ok := property["type"].(string)
			if !ok || !slices.Contains(supportedScalarTypes, propertyType) {
				return fmt.Errorf("input_schema.properties.%s.type must be one of %v", name, supportedScalarTypes)
			}

			if enumValues, ok := property["enum"]; ok {
				values, ok := enumValues.([]any)
				if !ok || len(values) == 0 {
					return fmt.Errorf("input_schema.properties.%s.enum must be a non-empty array", name)
				}
				for _, value := range values {
					if !valueMatchesType(propertyType, value) {
						return fmt.Errorf("input_schema.properties.%s.enum contains a value that does not match type %q", name, propertyType)
					}
				}
			}
		}
	}

	if requiredValue, ok := schema["required"]; ok {
		requiredFields, ok := requiredValue.([]any)
		if !ok {
			return fmt.Errorf("input_schema.required must be an array")
		}

		properties, _ := schema["properties"].(map[string]any)
		for _, fieldValue := range requiredFields {
			fieldName, ok := fieldValue.(string)
			if !ok || fieldName == "" {
				return fmt.Errorf("input_schema.required entries must be non-empty strings")
			}
			if properties != nil {
				if _, ok := properties[fieldName]; !ok {
					return fmt.Errorf("input_schema.required references unknown property %q", fieldName)
				}
			}
		}
	}

	return nil
}

func valueMatchesType(schemaType string, value any) bool {
	switch schemaType {
	case "string":
		_, ok := value.(string)
		return ok
	case "integer":
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		case float64:
			return math.Trunc(value.(float64)) == value.(float64)
		default:
			return false
		}
	case "number":
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			return true
		default:
			return false
		}
	case "boolean":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}
