package templates

import (
	"github.com/LederWorks/hippodamus/pkg/schema"
)

// ShapeTemplate handles template generation for shape elements
type ShapeTemplate struct{}

// NewShapeTemplate creates a new shape template generator
func NewShapeTemplate() *ShapeTemplate {
	return &ShapeTemplate{}
}

// Generate creates a schema.Element from shape parameters
func (t *ShapeTemplate) Generate(params map[string]interface{}) (*schema.Element, error) {
	label := getStringParam(params, "label", "Shape Element")
	shape := getStringParam(params, "shape", "rectangle")
	fillColor := getStringParam(params, "fillColor", "#E3F2FD")
	strokeColor := getStringParam(params, "strokeColor", "#1976D2")

	width := getFloatParam(params, "width", 120)
	height := getFloatParam(params, "height", 80)
	x := getFloatParam(params, "x", 100)
	y := getFloatParam(params, "y", 100)
	strokeWidth := getFloatParam(params, "strokeWidth", 2)
	fontSize := getIntParam(params, "fontSize", 14)

	fontStyle := getStringParam(params, "fontStyle", "normal")
	rounded := getBoolParam(params, "rounded", true)
	shadow := getBoolParam(params, "shadow", false)

	return &schema.Element{
		Type: schema.ElementTypeShape,
		Properties: schema.ElementProperties{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
			Label:  label,
			Shape:  shape,
		},
		Style: schema.Style{
			FillColor:   fillColor,
			StrokeColor: strokeColor,
			StrokeWidth: strokeWidth,
			FontSize:    fontSize,
			FontStyle:   fontStyle,
			Rounded:     rounded,
			Shadow:      shadow,
		},
	}, nil
}

// Helper functions for parameter extraction
func getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return defaultValue
}

func getIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return defaultValue
}

func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}
