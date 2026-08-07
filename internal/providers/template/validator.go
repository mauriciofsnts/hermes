package template

import (
	"regexp"
	"strings"
)

var templateVarRegex = regexp.MustCompile(`\{\{\.(\w+)}}`)

// ValidateTemplateStructure validates that the provided data contains the variables required by the template.
func ValidateTemplateStructure(templateContent []byte, requiredData map[string]any) error {
	matches := templateVarRegex.FindAllStringSubmatch(string(templateContent), -1)

	for _, match := range matches {
		fieldName := match[1]

		if isControlField(fieldName) {
			continue
		}

		if _, exists := requiredData[fieldName]; !exists {
			return NewTemplateMissingFieldError(fieldName)
		}
	}

	return nil
}

// ExtractTemplateVariables extracts all variables from a template.
func ExtractTemplateVariables(templateContent []byte) []string {
	matches := templateVarRegex.FindAllStringSubmatch(string(templateContent), -1)

	vars := make(map[string]bool)
	for _, match := range matches {
		fieldName := match[1]
		if isControlField(fieldName) {
			continue
		}
		vars[fieldName] = true
	}

	result := make([]string, 0, len(vars))
	for v := range vars {
		result = append(result, v)
	}

	return result
}

func isControlField(fieldName string) bool {
	return fieldName == "range" || fieldName == "if" || fieldName == "else" ||
		strings.HasPrefix(fieldName, "range") || strings.HasPrefix(fieldName, "if")
}

// TemplateMissingFieldError represents a missing template field error.
type TemplateMissingFieldError struct {
	FieldName string
}

func NewTemplateMissingFieldError(fieldName string) *TemplateMissingFieldError {
	return &TemplateMissingFieldError{FieldName: fieldName}
}

func (e *TemplateMissingFieldError) Error() string {
	return "template requires field '" + e.FieldName + "' but it was not provided"
}
