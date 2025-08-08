package templates

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

func TestTextTemplate_Generate(t *testing.T) {
	template := NewTextTemplate()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected *schema.Element
	}{
		{
			name: "basic text",
			params: map[string]interface{}{
				"label": "Hello World",
			},
			expected: &schema.Element{
				Type: schema.ElementTypeText,
				Properties: schema.ElementProperties{
					X:      0,
					Y:      0,
					Width:  100,
					Height: 30,
					Label:  "Hello World",
				},
				Style: schema.Style{
					FontSize:      12,
					FontFamily:    "Arial",
					FontColor:     "#000000",
					FontStyle:     "normal",
					TextAlign:     "center",
					VerticalAlign: "middle",
					FillColor:     "",
					StrokeColor:   "",
					StrokeWidth:   0,
				},
			},
		},
		{
			name: "styled text",
			params: map[string]interface{}{
				"label":         "Styled Text",
				"x":             100.0,
				"y":             50.0,
				"width":         200.0,
				"height":        40.0,
				"fontSize":      16.0,
				"fontFamily":    "Times",
				"fontColor":     "#FF0000",
				"fontStyle":     "bold",
				"textAlign":     "left",
				"verticalAlign": "top",
				"fillColor":     "#FFFFFF",
				"strokeColor":   "#000000",
				"strokeWidth":   2.0,
			},
			expected: &schema.Element{
				Type: schema.ElementTypeText,
				Properties: schema.ElementProperties{
					X:      100,
					Y:      50,
					Width:  200,
					Height: 40,
					Label:  "Styled Text",
				},
				Style: schema.Style{
					FontSize:      16,
					FontFamily:    "Times",
					FontColor:     "#FF0000",
					FontStyle:     "bold",
					TextAlign:     "left",
					VerticalAlign: "top",
					FillColor:     "#FFFFFF",
					StrokeColor:   "#000000",
					StrokeWidth:   2,
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

			// Test style
			if result.Style.FontSize != tt.expected.Style.FontSize {
				t.Errorf("Expected FontSize %d, got %d", tt.expected.Style.FontSize, result.Style.FontSize)
			}

			if result.Style.FontFamily != tt.expected.Style.FontFamily {
				t.Errorf("Expected FontFamily %s, got %s", tt.expected.Style.FontFamily, result.Style.FontFamily)
			}

			if result.Style.FontColor != tt.expected.Style.FontColor {
				t.Errorf("Expected FontColor %s, got %s", tt.expected.Style.FontColor, result.Style.FontColor)
			}

			if result.Style.FontStyle != tt.expected.Style.FontStyle {
				t.Errorf("Expected FontStyle %s, got %s", tt.expected.Style.FontStyle, result.Style.FontStyle)
			}

			if result.Style.TextAlign != tt.expected.Style.TextAlign {
				t.Errorf("Expected TextAlign %s, got %s", tt.expected.Style.TextAlign, result.Style.TextAlign)
			}

			if result.Style.VerticalAlign != tt.expected.Style.VerticalAlign {
				t.Errorf("Expected VerticalAlign %s, got %s", tt.expected.Style.VerticalAlign, result.Style.VerticalAlign)
			}

			if result.Style.FillColor != tt.expected.Style.FillColor {
				t.Errorf("Expected FillColor %s, got %s", tt.expected.Style.FillColor, result.Style.FillColor)
			}

			if result.Style.StrokeColor != tt.expected.Style.StrokeColor {
				t.Errorf("Expected StrokeColor %s, got %s", tt.expected.Style.StrokeColor, result.Style.StrokeColor)
			}

			if result.Style.StrokeWidth != tt.expected.Style.StrokeWidth {
				t.Errorf("Expected StrokeWidth %f, got %f", tt.expected.Style.StrokeWidth, result.Style.StrokeWidth)
			}
		})
	}
}

func TestTextTemplate_DefaultValues(t *testing.T) {
	template := NewTextTemplate()

	// Test with minimal parameters
	params := map[string]interface{}{
		"label": "Test",
	}

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

	if result.Properties.Width != 100 {
		t.Errorf("Expected default Width 100, got %f", result.Properties.Width)
	}

	if result.Properties.Height != 30 {
		t.Errorf("Expected default Height 30, got %f", result.Properties.Height)
	}

	if result.Style.FontSize != 12 {
		t.Errorf("Expected default FontSize 12, got %d", result.Style.FontSize)
	}

	if result.Style.FontFamily != "Arial" {
		t.Errorf("Expected default FontFamily 'Arial', got %s", result.Style.FontFamily)
	}

	if result.Style.FontColor != "#000000" {
		t.Errorf("Expected default FontColor '#000000', got %s", result.Style.FontColor)
	}

	if result.Style.FontStyle != "normal" {
		t.Errorf("Expected default FontStyle 'normal', got %s", result.Style.FontStyle)
	}

	if result.Style.TextAlign != "center" {
		t.Errorf("Expected default TextAlign 'center', got %s", result.Style.TextAlign)
	}

	if result.Style.VerticalAlign != "middle" {
		t.Errorf("Expected default VerticalAlign 'middle', got %s", result.Style.VerticalAlign)
	}
}
