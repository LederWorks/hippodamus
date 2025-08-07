package resources

import (
	"github.com/LederWorks/hippodamus/pkg/providers"
)

// SwimlaneResource defines the swimlane container resource
type SwimlaneResource struct{}

// NewSwimlaneResource creates a new swimlane resource instance
func NewSwimlaneResource() *SwimlaneResource {
	return &SwimlaneResource{}
}

// Definition returns the resource definition for swimlane elements
func (r *SwimlaneResource) Definition() providers.ResourceDefinition {
	return providers.ResourceDefinition{
		Type:        "swimlane",
		Name:        "Swimlane",
		Description: "Horizontal or vertical lane for organizing process flows",
		Category:    "container",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "Swimlane title/label",
					"default":     "",
				},
				"x": map[string]interface{}{
					"type":        "number",
					"description": "X position",
					"default":     0,
				},
				"y": map[string]interface{}{
					"type":        "number",
					"description": "Y position",
					"default":     0,
				},
				"width": map[string]interface{}{
					"type":        "number",
					"description": "Swimlane width",
					"default":     300,
					"minimum":     100,
				},
				"height": map[string]interface{}{
					"type":        "number",
					"description": "Swimlane height",
					"default":     200,
					"minimum":     50,
				},
				"orientation": map[string]interface{}{
					"type":        "string",
					"description": "Swimlane orientation",
					"default":     "horizontal",
					"enum":        []string{"horizontal", "vertical"},
				},
				"startSize": map[string]interface{}{
					"type":        "number",
					"description": "Size of the header area",
					"default":     30,
					"minimum":     20,
				},
				"fillColor": map[string]interface{}{
					"type":        "string",
					"description": "Background color",
					"default":     "#F8F9FA",
				},
				"strokeColor": map[string]interface{}{
					"type":        "string",
					"description": "Border color",
					"default":     "#6C757D",
				},
				"strokeWidth": map[string]interface{}{
					"type":        "number",
					"description": "Border width",
					"default":     1,
					"minimum":     0,
				},
				"collapsible": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the swimlane can be collapsed",
					"default":     true,
				},
				"collapsed": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the swimlane starts collapsed",
					"default":     false,
				},
				"fontSize": map[string]interface{}{
					"type":        "integer",
					"description": "Font size for the label",
					"default":     12,
					"minimum":     6,
					"maximum":     72,
				},
				"fontStyle": map[string]interface{}{
					"type":        "string",
					"description": "Font style for the label",
					"default":     "bold",
					"enum":        []string{"normal", "bold", "italic", "bold italic"},
				},
				"fontColor": map[string]interface{}{
					"type":        "string",
					"description": "Font color for the label",
					"default":     "#000000",
				},
				"childLayout": map[string]interface{}{
					"type":        "string",
					"description": "How children are laid out",
					"default":     "stackLayout",
					"enum":        []string{"stackLayout", "flowLayout", "freeLayout"},
				},
			},
			"required": []string{},
		},
		Examples: []providers.ResourceExample{
			{
				Name:        "Horizontal Swimlane",
				Description: "Basic horizontal swimlane for process flows",
				Config: map[string]interface{}{
					"label":       "Customer Service",
					"x":           50,
					"y":           100,
					"width":       400,
					"height":      150,
					"orientation": "horizontal",
					"startSize":   35,
					"fillColor":   "#E3F2FD",
					"strokeColor": "#1976D2",
				},
			},
			{
				Name:        "Vertical Swimlane",
				Description: "Vertical swimlane for role-based organization",
				Config: map[string]interface{}{
					"label":       "Development Team",
					"x":           200,
					"y":           50,
					"width":       150,
					"height":      300,
					"orientation": "vertical",
					"startSize":   40,
					"fillColor":   "#F3E5F5",
					"strokeColor": "#7B1FA2",
					"collapsible": true,
				},
			},
		},
	}
}

// Validate validates swimlane parameters
func (r *SwimlaneResource) Validate(params map[string]interface{}) error {
	// Validate orientation if provided
	if orientation, ok := params["orientation"]; ok {
		validOrientations := []string{"horizontal", "vertical"}
		if orient, ok := orientation.(string); ok {
			valid := false
			for _, validOrient := range validOrientations {
				if orient == validOrient {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "orientation",
					Message: "invalid orientation",
				}
			}
		}
	}

	// Validate child layout if provided
	if childLayout, ok := params["childLayout"]; ok {
		validLayouts := []string{"stackLayout", "flowLayout", "freeLayout"}
		if layout, ok := childLayout.(string); ok {
			valid := false
			for _, validLayout := range validLayouts {
				if layout == validLayout {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "childLayout",
					Message: "invalid child layout",
				}
			}
		}
	}

	// Validate font style if provided
	if fontStyle, ok := params["fontStyle"]; ok {
		validStyles := []string{"normal", "bold", "italic", "bold italic"}
		if style, ok := fontStyle.(string); ok {
			valid := false
			for _, validStyle := range validStyles {
				if style == validStyle {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "fontStyle",
					Message: "invalid font style",
				}
			}
		}
	}

	// Validate dimensions if provided
	if width, ok := params["width"]; ok {
		if w, ok := width.(float64); ok && w < 100 {
			return &providers.ValidationError{
				Field:   "width",
				Message: "width must be at least 100",
			}
		}
	}

	if height, ok := params["height"]; ok {
		if h, ok := height.(float64); ok && h < 50 {
			return &providers.ValidationError{
				Field:   "height",
				Message: "height must be at least 50",
			}
		}
	}

	// Validate start size if provided
	if startSize, ok := params["startSize"]; ok {
		if ss, ok := startSize.(float64); ok && ss < 20 {
			return &providers.ValidationError{
				Field:   "startSize",
				Message: "start size must be at least 20",
			}
		}
	}

	// Validate font size if provided
	if fontSize, ok := params["fontSize"]; ok {
		if fs, ok := fontSize.(float64); ok && (fs < 6 || fs > 72) {
			return &providers.ValidationError{
				Field:   "fontSize",
				Message: "font size must be between 6 and 72",
			}
		}
	}

	// Validate stroke width if provided
	if strokeWidth, ok := params["strokeWidth"]; ok {
		if sw, ok := strokeWidth.(float64); ok && sw < 0 {
			return &providers.ValidationError{
				Field:   "strokeWidth",
				Message: "stroke width cannot be negative",
			}
		}
	}

	return nil
}
