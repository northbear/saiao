package invoke

import (
	"fmt"
	"math"

	saiaoerrors "saiao/internal/errors"
)

func ValidateInput(schema map[string]any, input map[string]any) error {
	if len(schema) == 0 {
		return nil
	}

	requiredFields, _ := schema["required"].([]any)
	for _, fieldValue := range requiredFields {
		fieldName := fieldValue.(string)
		if _, ok := input[fieldName]; !ok {
			return fmt.Errorf("%w: missing required field %q", saiaoerrors.ErrInvalidInput, fieldName)
		}
	}

	properties, _ := schema["properties"].(map[string]any)
	for fieldName, value := range input {
		propertyValue, ok := properties[fieldName]
		if !ok {
			continue
		}

		property := propertyValue.(map[string]any)
		propertyType := property["type"].(string)
		if !matchesType(propertyType, value) {
			return fmt.Errorf("%w: field %q must be %s", saiaoerrors.ErrInvalidInput, fieldName, propertyType)
		}

		enumValues, ok := property["enum"].([]any)
		if ok && len(enumValues) > 0 {
			if !matchesEnum(enumValues, value) {
				return fmt.Errorf("%w: field %q must be one of the configured enum values", saiaoerrors.ErrInvalidInput, fieldName)
			}
		}
	}

	return nil
}

func matchesType(schemaType string, value any) bool {
	switch schemaType {
	case "string":
		_, ok := value.(string)
		return ok
	case "integer":
		number, ok := value.(float64)
		return ok && math.Trunc(number) == number
	case "number":
		_, ok := value.(float64)
		return ok
	case "boolean":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

func matchesEnum(values []any, candidate any) bool {
	for _, value := range values {
		if fmt.Sprint(value) == fmt.Sprint(candidate) {
			return true
		}
	}
	return false
}
