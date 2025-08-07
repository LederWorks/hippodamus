package resources

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/providers"
)

func TestGroupResource_Definition(t *testing.T) {
	resource := NewGroupResource()
	def := resource.Definition()

	// Test basic properties
	if def.Type != "group" {
		t.Errorf("Expected type 'group', got %s", def.Type)
	}

	if def.Name != "Group" {
		t.Errorf("Expected name 'Group', got %s", def.Name)
	}

	if def.Category != "container" {
		t.Errorf("Expected category 'container', got %s", def.Category)
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
	requiredFields := []string{"label", "x", "y", "width", "height", "fillColor", "strokeColor", "strokeWidth", "strokeStyle", "rounded", "collapsible", "collapsed"}
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

func TestGroupResource_Validate(t *testing.T) {
	resource := NewGroupResource()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name:        "valid empty group",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name: "valid group with all params",
			params: map[string]interface{}{
				"label":       "Test Group",
				"x":           100.0,
				"y":           50.0,
				"width":       200.0,
				"height":      150.0,
				"fillColor":   "#F0F0F0",
				"strokeColor": "#000000",
				"strokeWidth": 2.0,
				"strokeStyle": "solid",
				"rounded":     true,
				"collapsible": true,
				"collapsed":   false,
				"fontSize":    14.0,
				"fontStyle":   "bold",
				"fontColor":   "#333333",
			},
			expectError: false,
		},
		{
			name: "invalid stroke style",
			params: map[string]interface{}{
				"strokeStyle": "invalid",
			},
			expectError: true,
			errorField:  "strokeStyle",
		},
		{
			name: "invalid font style",
			params: map[string]interface{}{
				"fontStyle": "invalid",
			},
			expectError: true,
			errorField:  "fontStyle",
		},
		{
			name: "width too small",
			params: map[string]interface{}{
				"width": 30.0,
			},
			expectError: true,
			errorField:  "width",
		},
		{
			name: "height too small",
			params: map[string]interface{}{
				"height": 30.0,
			},
			expectError: true,
			errorField:  "height",
		},
		{
			name: "font size too small",
			params: map[string]interface{}{
				"fontSize": 3.0,
			},
			expectError: true,
			errorField:  "fontSize",
		},
		{
			name: "font size too large",
			params: map[string]interface{}{
				"fontSize": 100.0,
			},
			expectError: true,
			errorField:  "fontSize",
		},
		{
			name: "negative stroke width",
			params: map[string]interface{}{
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

func TestGroupResource_ValidOptions(t *testing.T) {
	resource := NewGroupResource()

	// Test all valid stroke styles
	validStrokeStyles := []string{"solid", "dashed", "dotted"}
	for _, style := range validStrokeStyles {
		params := map[string]interface{}{
			"strokeStyle": style,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid stroke style %s, got: %v", style, err)
		}
	}

	// Test all valid font styles
	validFontStyles := []string{"normal", "bold", "italic", "bold italic"}
	for _, style := range validFontStyles {
		params := map[string]interface{}{
			"fontStyle": style,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid font style %s, got: %v", style, err)
		}
	}
}
