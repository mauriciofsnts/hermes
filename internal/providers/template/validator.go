package template

import (
	"regexp"
	"strings"
)

// ValidateTemplateStructure valida se os dados possuem as chaves necessárias para o template
func ValidateTemplateStructure(templateContent []byte, requiredData map[string]any) error {
	// Extrair variáveis do template usando regex
	// Procura por padrões como {{.FieldName}}
	re := regexp.MustCompile(`\{\{\.(\w+)}}`)
	matches := re.FindAllStringSubmatch(string(templateContent), -1)

	// Verificar se todas as variáveis necessárias estão presentes
	for _, match := range matches {
		fieldName := match[1]

		// Ignorar campos padrão
		if fieldName == "range" || fieldName == "if" || fieldName == "else" {
			continue
		}

		if _, exists := requiredData[fieldName]; !exists {
			return NewTemplateMissingFieldError(fieldName)
		}
	}

	return nil
}

// ExtractTemplateVariables extrai todas as variáveis de um template
func ExtractTemplateVariables(templateContent []byte) []string {
	re := regexp.MustCompile(`\{\{\.(\w+)}}`)
	matches := re.FindAllStringSubmatch(string(templateContent), -1)

	// Usar map para evitar duplicatas
	vars := make(map[string]bool)
	for _, match := range matches {
		fieldName := match[1]
		// Ignorar campos padrão
		if !strings.HasPrefix(fieldName, "range") && !strings.HasPrefix(fieldName, "if") {
			vars[fieldName] = true
		}
	}

	// Converter map para slice
	result := make([]string, 0, len(vars))
	for v := range vars {
		result = append(result, v)
	}

	return result
}

// TemplateMissingFieldError representa um erro de campo faltante
type TemplateMissingFieldError struct {
	FieldName string
}

func NewTemplateMissingFieldError(fieldName string) *TemplateMissingFieldError {
	return &TemplateMissingFieldError{FieldName: fieldName}
}

func (e *TemplateMissingFieldError) Error() string {
	return "template requires field '" + e.FieldName + "' but it was not provided"
}
