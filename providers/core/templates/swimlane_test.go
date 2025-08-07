package templates

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

func TestSwimlaneTemplate_Generate(t *testing.T) {
	template := NewSwimlaneTemplate()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected *schema.Element
	}{
		{
			name:   "basic swimlane",
			params: map[string]interface{}{},
			expected: &schema.Element{
				Type: schema.ElementTypeSwimLane,
				Properties: schema.ElementProperties{
					X:           0,
					Y:           0,
					Width:       300,
					Height:      200,
					Label:       "",
					Collapsible: true,
					Collapsed:   false,
					Custom: map[string]interface{}{
						"orientation": "horizontal",
						"startSize":   30.0,
						"childLayout": "stackLayout",
						"horizontal":  true,
					},
				},
				Style: schema.Style{
					FillColor:   "#F8F9FA",
					StrokeColor: "#6C757D",
					StrokeWidth: 1,
					FontSize:    12,
					FontStyle:   "bold",
					FontColor:   "#000000",
				},
			},
		},
		{
			name: "vertical swimlane",
			params: map[string]interface{}{
				"label":       "Test Lane",
				"x":           100.0,
				"y":           50.0,
				"width":       150.0,
				"height":      400.0,
				"orientation": "vertical",
				"startSize":   40.0,
				"fillColor":   "#E0E0E0",
				"strokeColor": "#FF0000",
				"strokeWidth": 2.0,
				"collapsible": false,
				"collapsed":   false,
				"fontSize":    14.0,
				"fontStyle":   "italic",
				"fontColor":   "#333333",
				"childLayout": "flowLayout",
			},
			expected: &schema.Element{
				Type: schema.ElementTypeSwimLane,
				Properties: schema.ElementProperties{
					X:           100,
					Y:           50,
					Width:       150,
					Height:      400,
					Label:       "Test Lane",
					Collapsible: false,
					Collapsed:   false,
					Custom: map[string]interface{}{
						"orientation": "vertical",
						"startSize":   40.0,
						"childLayout": "flowLayout",
						"horizontal":  false,
					},
				},
				Style: schema.Style{
					FillColor:   "#E0E0E0",
					StrokeColor: "#FF0000",
					StrokeWidth: 2,
					FontSize:    14,
					FontStyle:   "italic",
					FontColor:   "#333333",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := template.Generate(tt.params)

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result, got nil")
				return
			}

			// Test element type
			if result.Type != tt.expected.Type {
				t.Errorf("Expected type %s, got %s", tt.expected.Type, result.Type)
			}

			// Test properties
			if result.Properties.X != tt.expected.Properties.X {
				t.Errorf("Expected X %f, got %f", tt.expected.Properties.X, result.Properties.X)
			}

			if result.Properties.Y != tt.expected.Properties.Y {
				t.Errorf("Expected Y %f, got %f", tt.expected.Properties.Y, result.Properties.Y)
			}

			if result.Properties.Width != tt.expected.Properties.Width {
				t.Errorf("Expected Width %f, got %f", tt.expected.Properties.Width, result.Properties.Width)
			}

			if result.Properties.Height != tt.expected.Properties.Height {
				t.Errorf("Expected Height %f, got %f", tt.expected.Properties.Height, result.Properties.Height)
			}

			if result.Properties.Label != tt.expected.Properties.Label {
				t.Errorf("Expected Label %s, got %s", tt.expected.Properties.Label, result.Properties.Label)
			}

			if result.Properties.Collapsible != tt.expected.Properties.Collapsible {
				t.Errorf("Expected Collapsible %t, got %t", tt.expected.Properties.Collapsible, result.Properties.Collapsible)
			}

			if result.Properties.Collapsed != tt.expected.Properties.Collapsed {
				t.Errorf("Expected Collapsed %t, got %t", tt.expected.Properties.Collapsed, result.Properties.Collapsed)
			}

			// Test custom properties
			if result.Properties.Custom != nil && tt.expected.Properties.Custom != nil {
				if orientation := result.Properties.Custom["orientation"]; orientation != tt.expected.Properties.Custom["orientation"] {
					t.Errorf("Expected orientation %v, got %v", tt.expected.Properties.Custom["orientation"], orientation)
				}

				if startSize := result.Properties.Custom["startSize"]; startSize != tt.expected.Properties.Custom["startSize"] {
					t.Errorf("Expected startSize %v, got %v", tt.expected.Properties.Custom["startSize"], startSize)
				}

				if childLayout := result.Properties.Custom["childLayout"]; childLayout != tt.expected.Properties.Custom["childLayout"] {
					t.Errorf("Expected childLayout %v, got %v", tt.expected.Properties.Custom["childLayout"], childLayout)
				}

				if horizontal := result.Properties.Custom["horizontal"]; horizontal != tt.expected.Properties.Custom["horizontal"] {
					t.Errorf("Expected horizontal %v, got %v", tt.expected.Properties.Custom["horizontal"], horizontal)
				}
			}

			// Test style
			if result.Style.FillColor != tt.expected.Style.FillColor {
				t.Errorf("Expected FillColor %s, got %s", tt.expected.Style.FillColor, result.Style.FillColor)
			}

			if result.Style.StrokeColor != tt.expected.Style.StrokeColor {
				t.Errorf("Expected StrokeColor %s, got %s", tt.expected.Style.StrokeColor, result.Style.StrokeColor)
			}

			if result.Style.StrokeWidth != tt.expected.Style.StrokeWidth {
				t.Errorf("Expected StrokeWidth %f, got %f", tt.expected.Style.StrokeWidth, result.Style.StrokeWidth)
			}

			if result.Style.FontSize != tt.expected.Style.FontSize {
				t.Errorf("Expected FontSize %d, got %d", tt.expected.Style.FontSize, result.Style.FontSize)
			}

			if result.Style.FontStyle != tt.expected.Style.FontStyle {
				t.Errorf("Expected FontStyle %s, got %s", tt.expected.Style.FontStyle, result.Style.FontStyle)
			}

			if result.Style.FontColor != tt.expected.Style.FontColor {
				t.Errorf("Expected FontColor %s, got %s", tt.expected.Style.FontColor, result.Style.FontColor)
			}
		})
	}
}

func TestSwimlaneTemplate_DefaultValues(t *testing.T) {
	template := NewSwimlaneTemplate()

	// Test with empty parameters
	params := map[string]interface{}{}

	result, err := template.Generate(params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify default values
	if result.Properties.X != 0 {
		t.Errorf("Expected default X 0, got %f", result.Properties.X)
	}

	if result.Properties.Y != 0 {
		t.Errorf("Expected default Y 0, got %f", result.Properties.Y)
	}

	if result.Properties.Width != 300 {
		t.Errorf("Expected default Width 300, got %f", result.Properties.Width)
	}

	if result.Properties.Height != 200 {
		t.Errorf("Expected default Height 200, got %f", result.Properties.Height)
	}

	if result.Properties.Label != "" {
		t.Errorf("Expected default Label '', got %s", result.Properties.Label)
	}

	if result.Properties.Collapsible != true {
		t.Errorf("Expected default Collapsible true, got %t", result.Properties.Collapsible)
	}

	if result.Properties.Collapsed != false {
		t.Errorf("Expected default Collapsed false, got %t", result.Properties.Collapsed)
	}

	if result.Style.FillColor != "#F8F9FA" {
		t.Errorf("Expected default FillColor '#F8F9FA', got %s", result.Style.FillColor)
	}

	if result.Style.StrokeColor != "#6C757D" {
		t.Errorf("Expected default StrokeColor '#6C757D', got %s", result.Style.StrokeColor)
	}

	if result.Style.StrokeWidth != 1 {
		t.Errorf("Expected default StrokeWidth 1, got %f", result.Style.StrokeWidth)
	}

	if result.Style.FontSize != 12 {
		t.Errorf("Expected default FontSize 12, got %d", result.Style.FontSize)
	}

	if result.Style.FontStyle != "bold" {
		t.Errorf("Expected default FontStyle 'bold', got %s", result.Style.FontStyle)
	}

	if result.Properties.Custom["orientation"] != "horizontal" {
		t.Errorf("Expected default orientation 'horizontal', got %v", result.Properties.Custom["orientation"])
	}

	if result.Properties.Custom["startSize"] != 30.0 {
		t.Errorf("Expected default startSize 30.0, got %v", result.Properties.Custom["startSize"])
	}

	if result.Properties.Custom["childLayout"] != "stackLayout" {
		t.Errorf("Expected default childLayout 'stackLayout', got %v", result.Properties.Custom["childLayout"])
	}

	if result.Properties.Custom["horizontal"] != true {
		t.Errorf("Expected default horizontal true, got %v", result.Properties.Custom["horizontal"])
	}
}

func TestSwimlaneTemplate_OrientationMapping(t *testing.T) {
	template := NewSwimlaneTemplate()

	tests := []struct {
		orientation string
		horizontal  bool
	}{
		{"horizontal", true},
		{"vertical", false},
	}

	for _, tt := range tests {
		t.Run("orientation_"+tt.orientation, func(t *testing.T) {
			params := map[string]interface{}{
				"orientation": tt.orientation,
			}

			result, err := template.Generate(params)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if result.Properties.Custom["horizontal"] != tt.horizontal {
				t.Errorf("Expected horizontal %t for orientation %s, got %v", tt.horizontal, tt.orientation, result.Properties.Custom["horizontal"])
			}
		})
	}
}
