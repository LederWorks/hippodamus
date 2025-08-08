package core

import (
	"testing"
)

func TestCoreProvider_Basic(t *testing.T) {
	provider := NewCoreProvider()

	if provider.Name() != "core" {
		t.Errorf("Expected provider name 'core', got '%s'", provider.Name())
	}

	if provider.Version() != "dev" {
		t.Errorf("Expected version 'dev', got '%s'", provider.Version())
	}

	// Test that we have some resources (at least shape and connector)
	resources := provider.Resources()
	if len(resources) < 2 {
		t.Errorf("Expected at least 2 resources, got %d", len(resources))
	}

	// Verify shape resource exists
	hasShape := false
	hasConnector := false
	for _, resource := range resources {
		if resource.Type == "shape" {
			hasShape = true
		}
		if resource.Type == "connector" {
			hasConnector = true
		}
	}

	if !hasShape {
		t.Error("Expected shape resource to be available")
	}

	if !hasConnector {
		t.Error("Expected connector resource to be available")
	}
}

func TestCoreProvider_ShapeValidation(t *testing.T) {
	provider := NewCoreProvider()

	// Test valid shape parameters
	validParams := map[string]interface{}{
		"label": "Test Shape",
		"shape": "rectangle",
	}

	err := provider.Validate("shape", validParams)
	if err != nil {
		t.Errorf("Valid shape parameters should not cause error: %v", err)
	}

	// Test invalid shape parameters
	invalidParams := map[string]interface{}{
		"shape": "rectangle",
		// missing required label
	}

	err = provider.Validate("shape", invalidParams)
	if err == nil {
		t.Error("Invalid shape parameters should cause validation error")
	}
}

func TestCoreProvider_ShapeTemplate(t *testing.T) {
	provider := NewCoreProvider()

	params := map[string]interface{}{
		"label":       "Test Shape",
		"shape":       "ellipse",
		"fillColor":   "#FF0000",
		"strokeColor": "#000000",
		"width":       200.0,
		"height":      150.0,
	}

	element, err := provider.GenerateTemplate("shape", params)
	if err != nil {
		t.Fatalf("GenerateTemplate() error = %v", err)
	}

	if element.Properties.Label != "Test Shape" {
		t.Errorf("Expected label 'Test Shape', got '%s'", element.Properties.Label)
	}

	if element.Properties.Shape != "ellipse" {
		t.Errorf("Expected shape 'ellipse', got '%s'", element.Properties.Shape)
	}

	if element.Style.FillColor != "#FF0000" {
		t.Errorf("Expected fill color '#FF0000', got '%s'", element.Style.FillColor)
	}
}

func TestCoreProvider_ConnectorValidation(t *testing.T) {
	provider := NewCoreProvider()

	// Test valid connector parameters
	validParams := map[string]interface{}{
		"source": "element1",
		"target": "element2",
	}

	err := provider.Validate("connector", validParams)
	if err != nil {
		t.Errorf("Valid connector parameters should not cause error: %v", err)
	}

	// Test missing source
	invalidParams := map[string]interface{}{
		"target": "element2",
	}

	err = provider.Validate("connector", invalidParams)
	if err == nil {
		t.Error("Missing source should cause validation error")
	}
}

func TestCoreProvider_UnsupportedResource(t *testing.T) {
	provider := NewCoreProvider()

	err := provider.Validate("unsupported", map[string]interface{}{})
	if err == nil {
		t.Error("Unsupported resource should cause error")
	}

	_, err = provider.GenerateTemplate("unsupported", map[string]interface{}{})
	if err == nil {
		t.Error("Unsupported resource template should cause error")
	}
}
