package resources

import (
	"github.com/LederWorks/hippodamus/pkg/providers"
)

// GroupResource defines the group container resource
type GroupResource struct{}

// NewGroupResource creates a new group resource instance
func NewGroupResource() *GroupResource {
	return &GroupResource{}
}

// Definition returns the resource definition for group elements
func (r *GroupResource) Definition() providers.ResourceDefinition {
	return providers.ResourceDefinition{
		Type:        "group",
		Name:        "Group",
		Description: "Container element that groups related elements together",
		Category:    "container",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "Group title/label",
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
					"description": "Group width",
					"default":     200,
					"minimum":     50,
				},
				"height": map[string]interface{}{
					"type":        "number",
					"description": "Group height",
					"default":     150,
					"minimum":     50,
				},
				"fillColor": map[string]interface{}{
					"type":        "string",
					"description": "Background color",
					"default":     "#F5F5F5",
				},
				"strokeColor": map[string]interface{}{
					"type":        "string",
					"description": "Border color",
					"default":     "#CCCCCC",
				},
				"strokeWidth": map[string]interface{}{
					"type":        "number",
					"description": "Border width",
					"default":     1,
					"minimum":     0,
				},
				"strokeStyle": map[string]interface{}{
					"type":        "string",
					"description": "Border style",
					"default":     "solid",
					"enum":        []string{"solid", "dashed", "dotted"},
				},
				"rounded": map[string]interface{}{
					"type":        "boolean",
					"description": "Rounded corners",
					"default":     false,
				},
				"collapsible": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the group can be collapsed",
					"default":     false,
				},
				"collapsed": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the group starts collapsed",
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
			},
			"required": []string{},
		},
		Examples: []providers.ResourceExample{
			{
				Name:        "Simple Group",
				Description: "Basic group container",
				Config: map[string]interface{}{
					"label":       "Process Group",
					"x":           100,
					"y":           100,
					"width":       250,
					"height":      200,
					"fillColor":   "#E8F5E8",
					"strokeColor": "#4CAF50",
				},
			},
			{
				Name:        "Collapsible Group",
				Description: "Group that can be collapsed",
				Config: map[string]interface{}{
					"label":       "Advanced Settings",
					"x":           300,
					"y":           150,
					"width":       200,
					"height":      150,
					"fillColor":   "#FFF3E0",
					"strokeColor": "#FF9800",
					"strokeWidth": 2,
					"rounded":     true,
					"collapsible": true,
				},
			},
		},
	}
}

// Validate validates group parameters
func (r *GroupResource) Validate(params map[string]interface{}) error {
	// Validate stroke style if provided
	if strokeStyle, ok := params["strokeStyle"]; ok {
		validStyles := []string{"solid", "dashed", "dotted"}
		if style, ok := strokeStyle.(string); ok {
			valid := false
			for _, validStyle := range validStyles {
				if style == validStyle {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "strokeStyle",
					Message: "invalid stroke style",
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
		if w, ok := width.(float64); ok && w < 50 {
			return &providers.ValidationError{
				Field:   "width",
				Message: "width must be at least 50",
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
