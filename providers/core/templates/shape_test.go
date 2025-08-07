package templates

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

func TestShapeTemplate_Generate(t *testing.T) {
	template := NewShapeTemplate()

	t.Run("default parameters", func(t *testing.T) {
		params := map[string]interface{}{
			"label": "Test Shape",
		}

		element, err := template.Generate(params)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		if element.Type != schema.ElementTypeShape {
			t.Errorf("Expected type %s, got %s", schema.ElementTypeShape, element.Type)
		}

		if element.Properties.Label != "Test Shape" {
			t.Errorf("Expected label 'Test Shape', got '%s'", element.Properties.Label)
		}

		if element.Properties.Shape != "rectangle" {
			t.Errorf("Expected default shape 'rectangle', got '%s'", element.Properties.Shape)
		}

		if element.Properties.Width != 120 {
			t.Errorf("Expected default width 120, got %f", element.Properties.Width)
		}

		if element.Style.FillColor != "#E3F2FD" {
			t.Errorf("Expected default fill color '#E3F2FD', got '%s'", element.Style.FillColor)
		}
	})

	t.Run("custom parameters", func(t *testing.T) {
		params := map[string]interface{}{
			"label":       "Custom Shape",
			"shape":       "ellipse",
			"width":       200.0,
			"height":      150.0,
			"x":           50.0,
			"y":           75.0,
			"fillColor":   "#FF0000",
			"strokeColor": "#000000",
			"strokeWidth": 3.0,
			"fontSize":    18,
			"fontStyle":   "bold",
			"rounded":     false,
			"shadow":      true,
		}

		element, err := template.Generate(params)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		if element.Properties.Label != "Custom Shape" {
			t.Errorf("Expected label 'Custom Shape', got '%s'", element.Properties.Label)
		}

		if element.Properties.Shape != "ellipse" {
			t.Errorf("Expected shape 'ellipse', got '%s'", element.Properties.Shape)
		}

		if element.Properties.Width != 200.0 {
			t.Errorf("Expected width 200.0, got %f", element.Properties.Width)
		}

		if element.Properties.Height != 150.0 {
			t.Errorf("Expected height 150.0, got %f", element.Properties.Height)
		}

		if element.Properties.X != 50.0 {
			t.Errorf("Expected x 50.0, got %f", element.Properties.X)
		}

		if element.Properties.Y != 75.0 {
			t.Errorf("Expected y 75.0, got %f", element.Properties.Y)
		}

		if element.Style.FillColor != "#FF0000" {
			t.Errorf("Expected fill color '#FF0000', got '%s'", element.Style.FillColor)
		}

		if element.Style.StrokeColor != "#000000" {
			t.Errorf("Expected stroke color '#000000', got '%s'", element.Style.StrokeColor)
		}

		if element.Style.StrokeWidth != 3.0 {
			t.Errorf("Expected stroke width 3.0, got %f", element.Style.StrokeWidth)
		}

		if element.Style.FontSize != 18 {
			t.Errorf("Expected font size 18, got %d", element.Style.FontSize)
		}

		if element.Style.FontStyle != "bold" {
			t.Errorf("Expected font style 'bold', got '%s'", element.Style.FontStyle)
		}

		if element.Style.Rounded != false {
			t.Errorf("Expected rounded false, got %v", element.Style.Rounded)
		}

		if element.Style.Shadow != true {
			t.Errorf("Expected shadow true, got %v", element.Style.Shadow)
		}
	})
}

func TestShapeTemplate_ParameterTypes(t *testing.T) {
	template := NewShapeTemplate()

	t.Run("mixed parameter types", func(t *testing.T) {
		params := map[string]interface{}{
			"label":       "Test",
			"width":       int(100),      // int instead of float64
			"height":      int64(80),     // int64 instead of float64
			"strokeWidth": float64(2.5),  // float64
			"fontSize":    float64(16.0), // float64 instead of int
			"rounded":     true,          // bool
		}

		element, err := template.Generate(params)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		if element.Properties.Width != 100.0 {
			t.Errorf("Expected width 100.0 from int, got %f", element.Properties.Width)
		}

		if element.Properties.Height != 80.0 {
			t.Errorf("Expected height 80.0 from int64, got %f", element.Properties.Height)
		}

		if element.Style.StrokeWidth != 2.5 {
			t.Errorf("Expected stroke width 2.5, got %f", element.Style.StrokeWidth)
		}

		if element.Style.FontSize != 16 {
			t.Errorf("Expected font size 16 from float64, got %d", element.Style.FontSize)
		}

		if element.Style.Rounded != true {
			t.Errorf("Expected rounded true, got %v", element.Style.Rounded)
		}
	})
}

func TestShapeTemplate_HelperFunctions(t *testing.T) {
	t.Run("getStringParam", func(t *testing.T) {
		params := map[string]interface{}{
			"existing": "value",
			"number":   123,
		}

		result := getStringParam(params, "existing", "default")
		if result != "value" {
			t.Errorf("Expected 'value', got '%s'", result)
		}

		result = getStringParam(params, "missing", "default")
		if result != "default" {
			t.Errorf("Expected 'default', got '%s'", result)
		}

		result = getStringParam(params, "number", "default")
		if result != "default" {
			t.Errorf("Expected 'default' for non-string, got '%s'", result)
		}
	})

	t.Run("getFloatParam", func(t *testing.T) {
		params := map[string]interface{}{
			"float":  3.14,
			"int":    42,
			"int64":  int64(100),
			"string": "not_a_number",
		}

		result := getFloatParam(params, "float", 0.0)
		if result != 3.14 {
			t.Errorf("Expected 3.14, got %f", result)
		}

		result = getFloatParam(params, "int", 0.0)
		if result != 42.0 {
			t.Errorf("Expected 42.0, got %f", result)
		}

		result = getFloatParam(params, "int64", 0.0)
		if result != 100.0 {
			t.Errorf("Expected 100.0, got %f", result)
		}

		result = getFloatParam(params, "missing", 5.0)
		if result != 5.0 {
			t.Errorf("Expected 5.0, got %f", result)
		}

		result = getFloatParam(params, "string", 5.0)
		if result != 5.0 {
			t.Errorf("Expected 5.0 for non-number, got %f", result)
		}
	})

	t.Run("getBoolParam", func(t *testing.T) {
		params := map[string]interface{}{
			"true":   true,
			"false":  false,
			"string": "not_a_bool",
		}

		result := getBoolParam(params, "true", false)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}

		result = getBoolParam(params, "false", true)
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}

		result = getBoolParam(params, "missing", true)
		if result != true {
			t.Errorf("Expected true (default), got %v", result)
		}

		result = getBoolParam(params, "string", true)
		if result != true {
			t.Errorf("Expected true (default) for non-bool, got %v", result)
		}
	})
}
