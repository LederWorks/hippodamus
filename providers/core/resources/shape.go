package resources

import (
	"github.com/LederWorks/hippodamus/pkg/providers"
)

// ShapeResource defines the shape element resource
type ShapeResource struct{}

// NewShapeResource creates a new shape resource instance
func NewShapeResource() *ShapeResource {
	return &ShapeResource{}
}

// Definition returns the resource definition for shape elements
func (r *ShapeResource) Definition() providers.ResourceDefinition {
	return providers.ResourceDefinition{
		Type:        "shape",
		Name:        "Shape Element",
		Description: "Basic shape element with customizable appearance and properties",
		Category:    "basic",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "Text label for the shape",
					"default":     "Shape Element",
				},
				"shape": map[string]interface{}{
					"type":        "string",
					"description": "Shape type (rectangle, ellipse, triangle, diamond, etc.)",
					"default":     "rectangle",
					"enum":        []string{"rectangle", "ellipse", "triangle", "diamond", "hexagon", "cloud", "cylinder"},
				},
				"width": map[string]interface{}{
					"type":        "number",
					"description": "Width of the shape",
					"default":     120,
					"minimum":     10,
				},
				"height": map[string]interface{}{
					"type":        "number",
					"description": "Height of the shape",
					"default":     80,
					"minimum":     10,
				},
				"x": map[string]interface{}{
					"type":        "number",
					"description": "X position",
					"default":     100,
				},
				"y": map[string]interface{}{
					"type":        "number",
					"description": "Y position",
					"default":     100,
				},
				"fillColor": map[string]interface{}{
					"type":        "string",
					"description": "Fill color",
					"default":     "#E3F2FD",
				},
				"strokeColor": map[string]interface{}{
					"type":        "string",
					"description": "Border color",
					"default":     "#1976D2",
				},
				"strokeWidth": map[string]interface{}{
					"type":        "number",
					"description": "Border width",
					"default":     2,
					"minimum":     0,
				},
				"fontSize": map[string]interface{}{
					"type":        "number",
					"description": "Font size",
					"default":     14,
					"minimum":     8,
				},
				"fontStyle": map[string]interface{}{
					"type":        "string",
					"description": "Font style",
					"default":     "normal",
					"enum":        []string{"normal", "bold", "italic", "bold italic"},
				},
				"rounded": map[string]interface{}{
					"type":        "boolean",
					"description": "Enable rounded corners",
					"default":     true,
				},
				"shadow": map[string]interface{}{
					"type":        "boolean",
					"description": "Enable shadow effect",
					"default":     false,
				},
			},
			"required": []string{"label"},
		},
		Examples: []providers.ResourceExample{
			{
				Name:        "Basic Rectangle",
				Description: "Simple rectangular shape",
				Config: map[string]interface{}{
					"label":       "Basic Shape",
					"shape":       "rectangle",
					"fillColor":   "#E3F2FD",
					"strokeColor": "#1976D2",
				},
			},
			{
				Name:        "Rounded Cloud",
				Description: "Cloud shape with rounded appearance",
				Config: map[string]interface{}{
					"label":     "Cloud Service",
					"shape":     "cloud",
					"fillColor": "#FFF3E0",
					"rounded":   true,
					"shadow":    true,
				},
			},
		},
	}
}

// Validate validates shape parameters
func (r *ShapeResource) Validate(params map[string]interface{}) error {
	// Check required parameters
	if label, ok := params["label"]; !ok || label == "" {
		return &providers.ValidationError{
			Field:   "label",
			Message: "label is required",
		}
	}

	// Validate shape type if provided
	if shape, ok := params["shape"]; ok {
		validShapes := []string{"rectangle", "ellipse", "triangle", "diamond", "hexagon", "cloud", "cylinder"}
		shapeStr, ok := shape.(string)
		if !ok {
			return &providers.ValidationError{
				Field:   "shape",
				Message: "shape must be a string",
			}
		}

		valid := false
		for _, validShape := range validShapes {
			if shapeStr == validShape {
				valid = true
				break
			}
		}
		if !valid {
			return &providers.ValidationError{
				Field:   "shape",
				Message: "invalid shape type",
			}
		}
	}

	// Validate dimensions if provided
	if width, ok := params["width"]; ok {
		if w, ok := width.(float64); ok && w < 10 {
			return &providers.ValidationError{
				Field:   "width",
				Message: "width must be at least 10",
			}
		}
	}

	if height, ok := params["height"]; ok {
		if h, ok := height.(float64); ok && h < 10 {
			return &providers.ValidationError{
				Field:   "height",
				Message: "height must be at least 10",
			}
		}
	}

	return nil
}
