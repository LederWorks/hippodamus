package resources

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/providers"
)

func TestSwimlaneResource_Definition(t *testing.T) {
	resource := NewSwimlaneResource()
	def := resource.Definition()

	// Test basic properties
	if def.Type != "swimlane" {
		t.Errorf("Expected type 'swimlane', got %s", def.Type)
	}

	if def.Name != "Swimlane" {
		t.Errorf("Expected name 'Swimlane', got %s", def.Name)
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
	requiredFields := []string{"label", "x", "y", "width", "height", "orientation", "startSize", "fillColor", "strokeColor", "strokeWidth", "collapsible", "collapsed", "childLayout"}
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

func TestSwimlaneResource_Validate(t *testing.T) {
	resource := NewSwimlaneResource()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name:        "valid empty swimlane",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name: "valid horizontal swimlane",
			params: map[string]interface{}{
				"label":       "Test Lane",
				"x":           100.0,
				"y":           50.0,
				"width":       400.0,
				"height":      150.0,
				"orientation": "horizontal",
				"startSize":   30.0,
				"fillColor":   "#F0F0F0",
				"strokeColor": "#000000",
				"strokeWidth": 1.0,
				"collapsible": true,
				"collapsed":   false,
				"fontSize":    12.0,
				"fontStyle":   "bold",
				"fontColor":   "#333333",
				"childLayout": "stackLayout",
			},
			expectError: false,
		},
		{
			name: "valid vertical swimlane",
			params: map[string]interface{}{
				"orientation": "vertical",
				"childLayout": "flowLayout",
			},
			expectError: false,
		},
		{
			name: "invalid orientation",
			params: map[string]interface{}{
				"orientation": "diagonal",
			},
			expectError: true,
			errorField:  "orientation",
		},
		{
			name: "invalid child layout",
			params: map[string]interface{}{
				"childLayout": "invalid",
			},
			expectError: true,
			errorField:  "childLayout",
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
				"width": 50.0,
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
			name: "start size too small",
			params: map[string]interface{}{
				"startSize": 10.0,
			},
			expectError: true,
			errorField:  "startSize",
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

func TestSwimlaneResource_ValidOptions(t *testing.T) {
	resource := NewSwimlaneResource()

	// Test all valid orientations
	validOrientations := []string{"horizontal", "vertical"}
	for _, orientation := range validOrientations {
		params := map[string]interface{}{
			"orientation": orientation,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid orientation %s, got: %v", orientation, err)
		}
	}

	// Test all valid child layouts
	validLayouts := []string{"stackLayout", "flowLayout", "freeLayout"}
	for _, layout := range validLayouts {
		params := map[string]interface{}{
			"childLayout": layout,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid child layout %s, got: %v", layout, err)
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
