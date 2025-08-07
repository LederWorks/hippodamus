package templates

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

func TestConnectorTemplate_Generate(t *testing.T) {
	template := NewConnectorTemplate()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected *schema.Element
	}{
		{
			name: "basic connector",
			params: map[string]interface{}{
				"source": "element1",
				"target": "element2",
			},
			expected: &schema.Element{
				Type: schema.ElementTypeConnector,
				Properties: schema.ElementProperties{
					Source:     "element1",
					Target:     "element2",
					SourcePort: "right",
					TargetPort: "left",
					Label:      "",
					Custom: map[string]interface{}{
						"strokeStyle": "solid",
						"arrow":       "target",
					},
				},
				Style: schema.Style{
					StrokeColor: "#424242",
					StrokeWidth: 2,
				},
			},
		},
		{
			name: "connector with all params",
			params: map[string]interface{}{
				"source":      "start",
				"target":      "end",
				"sourcePort":  "top",
				"targetPort":  "bottom",
				"label":       "Data Flow",
				"strokeColor": "#FF0000",
				"strokeWidth": 3.0,
				"strokeStyle": "dashed",
				"arrow":       "both",
			},
			expected: &schema.Element{
				Type: schema.ElementTypeConnector,
				Properties: schema.ElementProperties{
					Source:     "start",
					Target:     "end",
					SourcePort: "top",
					TargetPort: "bottom",
					Label:      "Data Flow",
					Custom: map[string]interface{}{
						"strokeStyle": "dashed",
						"arrow":       "both",
					},
				},
				Style: schema.Style{
					StrokeColor: "#FF0000",
					StrokeWidth: 3,
				},
			},
		},
		{
			name: "connector with dotted style",
			params: map[string]interface{}{
				"source":      "node1",
				"target":      "node2",
				"strokeStyle": "dotted",
				"arrow":       "none",
			},
			expected: &schema.Element{
				Type: schema.ElementTypeConnector,
				Properties: schema.ElementProperties{
					Source:     "node1",
					Target:     "node2",
					SourcePort: "right",
					TargetPort: "left",
					Label:      "",
					Custom: map[string]interface{}{
						"strokeStyle": "dotted",
						"arrow":       "none",
					},
				},
				Style: schema.Style{
					StrokeColor: "#424242",
					StrokeWidth: 2,
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
			if result.Properties.Source != tt.expected.Properties.Source {
				t.Errorf("Expected source %s, got %s", tt.expected.Properties.Source, result.Properties.Source)
			}

			if result.Properties.Target != tt.expected.Properties.Target {
				t.Errorf("Expected target %s, got %s", tt.expected.Properties.Target, result.Properties.Target)
			}

			if result.Properties.SourcePort != tt.expected.Properties.SourcePort {
				t.Errorf("Expected sourcePort %s, got %s", tt.expected.Properties.SourcePort, result.Properties.SourcePort)
			}

			if result.Properties.TargetPort != tt.expected.Properties.TargetPort {
				t.Errorf("Expected targetPort %s, got %s", tt.expected.Properties.TargetPort, result.Properties.TargetPort)
			}

			if result.Properties.Label != tt.expected.Properties.Label {
				t.Errorf("Expected label %s, got %s", tt.expected.Properties.Label, result.Properties.Label)
			}

			// Test custom properties
			if result.Properties.Custom != nil && tt.expected.Properties.Custom != nil {
				if strokeStyle := result.Properties.Custom["strokeStyle"]; strokeStyle != tt.expected.Properties.Custom["strokeStyle"] {
					t.Errorf("Expected strokeStyle %v, got %v", tt.expected.Properties.Custom["strokeStyle"], strokeStyle)
				}

				if arrow := result.Properties.Custom["arrow"]; arrow != tt.expected.Properties.Custom["arrow"] {
					t.Errorf("Expected arrow %v, got %v", tt.expected.Properties.Custom["arrow"], arrow)
				}
			}

			// Test style
			if result.Style.StrokeColor != tt.expected.Style.StrokeColor {
				t.Errorf("Expected strokeColor %s, got %s", tt.expected.Style.StrokeColor, result.Style.StrokeColor)
			}

			if result.Style.StrokeWidth != tt.expected.Style.StrokeWidth {
				t.Errorf("Expected strokeWidth %f, got %f", tt.expected.Style.StrokeWidth, result.Style.StrokeWidth)
			}
		})
	}
}

func TestConnectorTemplate_DefaultValues(t *testing.T) {
	template := NewConnectorTemplate()

	// Test with minimal parameters
	params := map[string]interface{}{
		"source": "a",
		"target": "b",
	}

	result, err := template.Generate(params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify default values
	if result.Properties.SourcePort != "right" {
		t.Errorf("Expected default sourcePort 'right', got %s", result.Properties.SourcePort)
	}

	if result.Properties.TargetPort != "left" {
		t.Errorf("Expected default targetPort 'left', got %s", result.Properties.TargetPort)
	}

	if result.Properties.Label != "" {
		t.Errorf("Expected default label '', got %s", result.Properties.Label)
	}

	if result.Style.StrokeColor != "#424242" {
		t.Errorf("Expected default strokeColor '#424242', got %s", result.Style.StrokeColor)
	}

	if result.Style.StrokeWidth != 2 {
		t.Errorf("Expected default strokeWidth 2, got %f", result.Style.StrokeWidth)
	}

	if result.Properties.Custom["strokeStyle"] != "solid" {
		t.Errorf("Expected default strokeStyle 'solid', got %v", result.Properties.Custom["strokeStyle"])
	}

	if result.Properties.Custom["arrow"] != "target" {
		t.Errorf("Expected default arrow 'target', got %v", result.Properties.Custom["arrow"])
	}
}

func TestConnectorTemplate_PortValues(t *testing.T) {
	template := NewConnectorTemplate()

	ports := []string{"top", "right", "bottom", "left", "center"}

	for _, port := range ports {
		t.Run("port_"+port, func(t *testing.T) {
			params := map[string]interface{}{
				"source":     "a",
				"target":     "b",
				"sourcePort": port,
				"targetPort": port,
			}

			result, err := template.Generate(params)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if result.Properties.SourcePort != port {
				t.Errorf("Expected sourcePort %s, got %s", port, result.Properties.SourcePort)
			}

			if result.Properties.TargetPort != port {
				t.Errorf("Expected targetPort %s, got %s", port, result.Properties.TargetPort)
			}
		})
	}
}

func TestConnectorTemplate_StrokeStyles(t *testing.T) {
	template := NewConnectorTemplate()

	styles := []string{"solid", "dashed", "dotted"}

	for _, style := range styles {
		t.Run("style_"+style, func(t *testing.T) {
			params := map[string]interface{}{
				"source":      "a",
				"target":      "b",
				"strokeStyle": style,
			}

			result, err := template.Generate(params)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if result.Properties.Custom["strokeStyle"] != style {
				t.Errorf("Expected strokeStyle %s, got %v", style, result.Properties.Custom["strokeStyle"])
			}
		})
	}
}

func TestConnectorTemplate_ArrowStyles(t *testing.T) {
	template := NewConnectorTemplate()

	arrows := []string{"none", "source", "target", "both"}

	for _, arrow := range arrows {
		t.Run("arrow_"+arrow, func(t *testing.T) {
			params := map[string]interface{}{
				"source": "a",
				"target": "b",
				"arrow":  arrow,
			}

			result, err := template.Generate(params)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if result.Properties.Custom["arrow"] != arrow {
				t.Errorf("Expected arrow %s, got %v", arrow, result.Properties.Custom["arrow"])
			}
		})
	}
}
