package resources

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/providers"
)

func TestConnectorResource_Definition(t *testing.T) {
	resource := NewConnectorResource()
	def := resource.Definition()

	// Test basic properties
	if def.Type != "connector" {
		t.Errorf("Expected type 'connector', got %s", def.Type)
	}

	if def.Name != "Connector" {
		t.Errorf("Expected name 'Connector', got %s", def.Name)
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
	requiredFields := []string{"source", "target", "sourcePort", "targetPort", "label", "strokeColor", "strokeWidth", "strokeStyle", "arrow"}
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

func TestConnectorResource_Validate(t *testing.T) {
	resource := NewConnectorResource()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid basic connector",
			params: map[string]interface{}{
				"source": "element1",
				"target": "element2",
			},
			expectError: false,
		},
		{
			name: "valid connector with all params",
			params: map[string]interface{}{
				"source":      "element1",
				"target":      "element2",
				"sourcePort":  "right",
				"targetPort":  "left",
				"label":       "Connection",
				"strokeColor": "#FF0000",
				"strokeWidth": 3.0,
				"strokeStyle": "dashed",
				"arrow":       "both",
			},
			expectError: false,
		},
		{
			name: "missing source",
			params: map[string]interface{}{
				"target": "element2",
			},
			expectError: true,
			errorField:  "source",
		},
		{
			name: "missing target",
			params: map[string]interface{}{
				"source": "element1",
			},
			expectError: true,
			errorField:  "target",
		},
		{
			name: "empty source",
			params: map[string]interface{}{
				"source": "",
				"target": "element2",
			},
			expectError: true,
			errorField:  "source",
		},
		{
			name: "empty target",
			params: map[string]interface{}{
				"source": "element1",
				"target": "",
			},
			expectError: true,
			errorField:  "target",
		},
		{
			name: "invalid source port",
			params: map[string]interface{}{
				"source":     "element1",
				"target":     "element2",
				"sourcePort": "invalid",
			},
			expectError: true,
			errorField:  "sourcePort",
		},
		{
			name: "invalid target port",
			params: map[string]interface{}{
				"source":     "element1",
				"target":     "element2",
				"targetPort": "invalid",
			},
			expectError: true,
			errorField:  "targetPort",
		},
		{
			name: "invalid stroke style",
			params: map[string]interface{}{
				"source":      "element1",
				"target":      "element2",
				"strokeStyle": "invalid",
			},
			expectError: true,
			errorField:  "strokeStyle",
		},
		{
			name: "invalid arrow style",
			params: map[string]interface{}{
				"source": "element1",
				"target": "element2",
				"arrow":  "invalid",
			},
			expectError: true,
			errorField:  "arrow",
		},
		{
			name: "invalid stroke width",
			params: map[string]interface{}{
				"source":      "element1",
				"target":      "element2",
				"strokeWidth": 0.5,
			},
			expectError: true,
			errorField:  "strokeWidth",
		},
		{
			name: "valid ports",
			params: map[string]interface{}{
				"source":     "element1",
				"target":     "element2",
				"sourcePort": "top",
				"targetPort": "bottom",
			},
			expectError: false,
		},
		{
			name: "valid stroke styles",
			params: map[string]interface{}{
				"source":      "element1",
				"target":      "element2",
				"strokeStyle": "dotted",
			},
			expectError: false,
		},
		{
			name: "valid arrow styles",
			params: map[string]interface{}{
				"source": "element1",
				"target": "element2",
				"arrow":  "source",
			},
			expectError: false,
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

func TestConnectorResource_ValidateEdgeCases(t *testing.T) {
	resource := NewConnectorResource()

	// Test minimum stroke width
	params := map[string]interface{}{
		"source":      "element1",
		"target":      "element2",
		"strokeWidth": 1.0,
	}
	if err := resource.Validate(params); err != nil {
		t.Errorf("Expected no error for minimum stroke width, got: %v", err)
	}

	// Test all valid port values
	validPorts := []string{"top", "right", "bottom", "left", "center"}
	for _, port := range validPorts {
		params := map[string]interface{}{
			"source":     "element1",
			"target":     "element2",
			"sourcePort": port,
			"targetPort": port,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid port %s, got: %v", port, err)
		}
	}

	// Test all valid stroke styles
	validStyles := []string{"solid", "dashed", "dotted"}
	for _, style := range validStyles {
		params := map[string]interface{}{
			"source":      "element1",
			"target":      "element2",
			"strokeStyle": style,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid style %s, got: %v", style, err)
		}
	}

	// Test all valid arrow styles
	validArrows := []string{"none", "source", "target", "both"}
	for _, arrow := range validArrows {
		params := map[string]interface{}{
			"source": "element1",
			"target": "element2",
			"arrow":  arrow,
		}
		if err := resource.Validate(params); err != nil {
			t.Errorf("Expected no error for valid arrow %s, got: %v", arrow, err)
		}
	}
}
