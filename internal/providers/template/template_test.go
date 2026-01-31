package template

import (
	"testing"
)

func TestTemplateServiceCaching(t *testing.T) {
	service := NewTemplateService().(*TemplateService)

	// Criar um template de teste
	templateName := "test.html"
	content := []byte("<html><body>Hello {{.Name}}</body></html>")

	// Criar template
	err := service.Create(templateName, content)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}
	defer service.Delete(templateName)

	// Primeira renderização (deve ser carregada do disco)
	data := map[string]any{"Name": "World"}
	result1, err := service.ParseHtmlTemplate(templateName, data)
	if err != nil {
		t.Fatalf("Failed to parse template first time: %v", err)
	}

	// Verificar se foi armazenado em cache
	if len(service.cache) != 1 {
		t.Errorf("Expected 1 template in cache, got %d", len(service.cache))
	}

	// Segunda renderização (deve usar cache)
	result2, err := service.ParseHtmlTemplate(templateName, data)
	if err != nil {
		t.Fatalf("Failed to parse template second time: %v", err)
	}

	// Resultados devem ser iguais
	if result1.String() != result2.String() {
		t.Error("Results from cached and non-cached should be equal")
	}

	// Limpar cache
	err = service.ClearCache()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	if len(service.cache) != 0 {
		t.Errorf("Expected cache to be empty after clear, got %d items", len(service.cache))
	}
}

func TestTemplateCacheInvalidation(t *testing.T) {
	service := NewTemplateService().(*TemplateService)

	templateName := "test_invalidation.html"
	content1 := []byte("<html><body>Version 1: {{.Name}}</body></html>")

	// Criar template
	err := service.Create(templateName, content1)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}
	defer service.Delete(templateName)

	// Renderizar para cachear
	data := map[string]any{"Name": "Test"}
	service.ParseHtmlTemplate(templateName, data)

	if len(service.cache) != 1 {
		t.Errorf("Expected 1 item in cache after first parse, got %d", len(service.cache))
	}

	// Atualizar template
	content2 := []byte("<html><body>Version 2: {{.Name}}</body></html>")
	err = service.Create(templateName, content2)
	if err != nil {
		t.Fatalf("Failed to update template: %v", err)
	}

	// Cache deve ter sido invalidado
	if len(service.cache) != 0 {
		t.Error("Cache should have been invalidated after Create")
	}
}

func TestTemplateExists(t *testing.T) {
	service := NewTemplateService()

	templateName := "test_exists.html"
	content := []byte("<html><body>Test</body></html>")

	// Não deve existir antes
	if service.Exists(templateName) {
		t.Error("Template should not exist before creation")
	}

	// Criar
	err := service.Create(templateName, content)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}
	defer service.Delete(templateName)

	// Deve existir agora
	if !service.Exists(templateName) {
		t.Error("Template should exist after creation")
	}
}
