package resources

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/providers"
)

func TestTextResource_Definition(t *testing.T) {
	resource := NewTextResource()
	def := resource.Definition()

	// Test basic properties
	if def.Type != "text" {
		t.Errorf("Expected type 'text', got %s", def.Type)
	}

	if def.Name != "Text" {
		t.Errorf("Expected name 'Text', got %s", def.Name)
	}

	if def.Category != "basic" {
		t.Errorf("Expected category 'basic', got %s", def.Category)
	}

	// Test schema structure
	schema := def.Schema
	if schema == nil {
		t.Fatal("Schema should not be nil")
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema properties should be a map")
	}

	// Test required fields are defined
	requiredFields := []string{"label", "x", "y", "width", "height", "fontSize", "fontFamily", "fontColor", "fontStyle", "textAlign", "verticalAlign"}
	for _, field := range requiredFields {
		if _, exists := properties[field]; !exists {
			t.Errorf("Required field %s missing from schema", field)
		}
	}

	// Test examples
	if len(def.Examples) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(def.Examples))
	}
}

func TestTextResource_Validate(t *testing.T) {
	resource := NewTextResource()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid basic text",
			params: map[string]interface{}{
				"label": "Hello World",
			},
			expectError: false,
		},
		{
			name: "valid text with all params",
			params: map[string]interface{}{
				"label":         "Styled Text",
				"x":             100.0,
				"y":             50.0,
				"width":         150.0,
				"height":        40.0,
				"fontSize":      16.0,
				"fontFamily":    "Arial",
				"fontColor":     "#FF0000",
				"fontStyle":     "bold",
				"textAlign":     "center",
				"verticalAlign": "middle",
				"fillColor":     "#FFFFFF",
				"strokeColor":   "#000000",
				"strokeWidth":   1.0,
			},
			expectError: false,
		},
		{
			name:        "missing label",
			params:      map[string]interface{}{},
			expectError: true,
			errorField:  "label",
		},
		{
			name: "empty label",
			params: map[string]interface{}{
				"label": "",
			},
			expectError: true,
			errorField:  "label",
		},
		{
			name: "invalid font style",
			params: map[string]interface{}{
				"label":     "Test",
				"fontStyle": "invalid",
			},
			expectError: true,
			errorField:  "fontStyle",
		},
		{
			name: "invalid text align",
			params: map[string]interface{}{
				"label":     "Test",
				"textAlign": "invalid",
			},
			expectError: true,
			errorField:  "textAlign",
		},
		{
			name: "invalid vertical align",
			params: map[string]interface{}{
				"label":         "Test",
				"verticalAlign": "invalid",
			},
			expectError: true,
			errorField:  "verticalAlign",
		},
		{
			name: "width too small",
			params: map[string]interface{}{
				"label": "Test",
				"width": 5.0,
			},
			expectError: true,
			errorField:  "width",
		},
		{
			name: "height too small",
			params: map[string]interface{}{
				"label":  "Test",
				"height": 5.0,
			},
			expectError: true,
			errorField:  "height",
		},
		{
			name: "font size too small",
			params: map[string]interface{}{
				"label":    "Test",
				"fontSize": 3.0,
			},
			expectError: true,
			errorField:  "fontSize",
		},
		{
			name: "font size too large",
			params: map[string]interface{}{
				"label":    "Test",
				"fontSize": 100.0,
			},
			expectError: true,
			errorField:  "fontSize",
		},
		{
			name: "negative stroke width",
			params: map[string]interface{}{
				"label":       "Test",
				"strokeWidth": -1.0,
			},
			expectError: true,
			errorField:  "strokeWidth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resource.Validate(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error, but got none")
					return
				}

				if validationErr, ok := err.(*providers.ValidationError); ok {
					if validationErr.Field != tt.errorField {
						t.Errorf("Expected error field %s, got %s", tt.errorField, validationErr.Field)
					}
				} else {
					t.Errorf("Expected ValidationError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, but got: %v", err)
				}
			}
		})
	}
}

func TestTextResource_ValidOptions(t *testing.T) {
	resource := NewTextResource()

	// Test all valid font styles
	validFontStyles := []string{"normal", "bold", "italic", "bold italic"}
	for _, style := range validFontStyles {
		params := map[string]interface{}{
			"label":     "Test",
			"fontStyle": style,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid font style %s, got: %v", style, err)
		}
	}

	// Test all valid text alignments
	validTextAligns := []string{"left", "center", "right"}
	for _, align := range validTextAligns {
		params := map[string]interface{}{
			"label":     "Test",
			"textAlign": align,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid text align %s, got: %v", align, err)
		}
	}

	// Test all valid vertical alignments
	validVerticalAligns := []string{"top", "middle", "bottom"}
	for _, align := range validVerticalAligns {
		params := map[string]interface{}{
			"label":         "Test",
			"verticalAlign": align,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid vertical align %s, got: %v", align, err)
		}
	}
}
