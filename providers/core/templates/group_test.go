package templates

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

func TestGroupTemplate_Generate(t *testing.T) {
	template := NewGroupTemplate()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected *schema.Element
	}{
		{
			name:   "basic group",
			params: map[string]interface{}{},
			expected: &schema.Element{
				Type: schema.ElementTypeGroup,
				Properties: schema.ElementProperties{
					X:           0,
					Y:           0,
					Width:       200,
					Height:      150,
					Label:       "",
					Collapsible: false,
					Collapsed:   false,
				},
				Style: schema.Style{
					FillColor:   "#F5F5F5",
					StrokeColor: "#CCCCCC",
					StrokeWidth: 1,
					Rounded:     false,
					FontSize:    12,
					FontStyle:   "bold",
					FontColor:   "#000000",
					Custom: map[string]string{
						"strokeStyle": "solid",
					},
				},
			},
		},
		{
			name: "styled group",
			params: map[string]interface{}{
				"label":       "Test Group",
				"x":           100.0,
				"y":           50.0,
				"width":       300.0,
				"height":      200.0,
				"fillColor":   "#E0E0E0",
				"strokeColor": "#FF0000",
				"strokeWidth": 2.0,
				"strokeStyle": "dashed",
				"rounded":     true,
				"collapsible": true,
				"collapsed":   false,
				"fontSize":    14.0,
				"fontStyle":   "italic",
				"fontColor":   "#333333",
			},
			expected: &schema.Element{
				Type: schema.ElementTypeGroup,
				Properties: schema.ElementProperties{
					X:           100,
					Y:           50,
					Width:       300,
					Height:      200,
					Label:       "Test Group",
					Collapsible: true,
					Collapsed:   false,
				},
				Style: schema.Style{
					FillColor:   "#E0E0E0",
					StrokeColor: "#FF0000",
					StrokeWidth: 2,
					Rounded:     true,
					FontSize:    14,
					FontStyle:   "italic",
					FontColor:   "#333333",
					Custom: map[string]string{
						"strokeStyle": "dashed",
					},
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

			if result.Style.Rounded != tt.expected.Style.Rounded {
				t.Errorf("Expected Rounded %t, got %t", tt.expected.Style.Rounded, result.Style.Rounded)
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

			// Test custom properties
			if result.Style.Custom != nil && tt.expected.Style.Custom != nil {
				if strokeStyle := result.Style.Custom["strokeStyle"]; strokeStyle != tt.expected.Style.Custom["strokeStyle"] {
					t.Errorf("Expected strokeStyle %v, got %v", tt.expected.Style.Custom["strokeStyle"], strokeStyle)
				}
			}
		})
	}
}

func TestGroupTemplate_DefaultValues(t *testing.T) {
	template := NewGroupTemplate()

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

	if result.Properties.Width != 200 {
		t.Errorf("Expected default Width 200, got %f", result.Properties.Width)
	}

	if result.Properties.Height != 150 {
		t.Errorf("Expected default Height 150, got %f", result.Properties.Height)
	}

	if result.Properties.Label != "" {
		t.Errorf("Expected default Label '', got %s", result.Properties.Label)
	}

	if result.Properties.Collapsible != false {
		t.Errorf("Expected default Collapsible false, got %t", result.Properties.Collapsible)
	}

	if result.Properties.Collapsed != false {
		t.Errorf("Expected default Collapsed false, got %t", result.Properties.Collapsed)
	}

	if result.Style.FillColor != "#F5F5F5" {
		t.Errorf("Expected default FillColor '#F5F5F5', got %s", result.Style.FillColor)
	}

	if result.Style.StrokeColor != "#CCCCCC" {
		t.Errorf("Expected default StrokeColor '#CCCCCC', got %s", result.Style.StrokeColor)
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

	if result.Style.Custom["strokeStyle"] != "solid" {
		t.Errorf("Expected default strokeStyle 'solid', got %v", result.Style.Custom["strokeStyle"])
	}
}
