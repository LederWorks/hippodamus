package resources

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/providers"
)

func TestShapeResource_Definition(t *testing.T) {
	resource := NewShapeResource()
	def := resource.Definition()

	if def.Type != "shape" {
		t.Errorf("Expected type 'shape', got '%s'", def.Type)
	}

	if def.Name != "Shape Element" {
		t.Errorf("Expected name 'Shape Element', got '%s'", def.Name)
	}

	if def.Category != "basic" {
		t.Errorf("Expected category 'basic', got '%s'", def.Category)
	}

	// Test schema structure
	schema := def.Schema
	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected schema properties to be a map")
	}

	// Test required label property
	if _, hasLabel := properties["label"]; !hasLabel {
		t.Error("Expected schema to have label property")
	}

	// Test examples
	if len(def.Examples) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(def.Examples))
	}
}

func TestShapeResource_Validate(t *testing.T) {
	resource := NewShapeResource()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		errType string
	}{
		{
			name: "valid shape",
			params: map[string]interface{}{
				"label": "Test Shape",
				"shape": "rectangle",
			},
			wantErr: false,
		},
		{
			name: "missing label",
			params: map[string]interface{}{
				"shape": "rectangle",
			},
			wantErr: true,
			errType: "ValidationError",
		},
		{
			name: "empty label",
			params: map[string]interface{}{
				"label": "",
				"shape": "rectangle",
			},
			wantErr: true,
			errType: "ValidationError",
		},
		{
			name: "invalid shape type",
			params: map[string]interface{}{
				"label": "Test",
				"shape": "invalid_shape",
			},
			wantErr: true,
			errType: "ValidationError",
		},
		{
			name: "shape not string",
			params: map[string]interface{}{
				"label": "Test",
				"shape": 123,
			},
			wantErr: true,
			errType: "ValidationError",
		},
		{
			name: "width too small",
			params: map[string]interface{}{
				"label": "Test",
				"width": 5.0,
			},
			wantErr: true,
			errType: "ValidationError",
		},
		{
			name: "height too small",
			params: map[string]interface{}{
				"label":  "Test",
				"height": 5.0,
			},
			wantErr: true,
			errType: "ValidationError",
		},
		{
			name: "valid with all parameters",
			params: map[string]interface{}{
				"label":       "Complete Shape",
				"shape":       "ellipse",
				"width":       150.0,
				"height":      100.0,
				"x":           50.0,
				"y":           75.0,
				"fillColor":   "#FF0000",
				"strokeColor": "#000000",
				"strokeWidth": 3.0,
				"fontSize":    16.0,
				"fontStyle":   "bold",
				"rounded":     true,
				"shadow":      false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resource.Validate(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errType != "" {
				switch tt.errType {
				case "ValidationError":
					if _, ok := err.(*providers.ValidationError); !ok {
						t.Errorf("Expected ValidationError, got %T", err)
					}
				}
			}
		})
	}
}

func TestShapeResource_ValidShapes(t *testing.T) {
	resource := NewShapeResource()

	validShapes := []string{"rectangle", "ellipse", "triangle", "diamond", "hexagon", "cloud", "cylinder"}

	for _, shape := range validShapes {
		t.Run("shape_"+shape, func(t *testing.T) {
			params := map[string]interface{}{
				"label": "Test",
				"shape": shape,
			}

			err := resource.Validate(params)
			if err != nil {
				t.Errorf("Valid shape '%s' should not cause validation error: %v", shape, err)
			}
		})
	}
}
