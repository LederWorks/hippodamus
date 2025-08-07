package resources

import (
	"github.com/LederWorks/hippodamus/pkg/providers"
)

// TextResource defines the text element resource
type TextResource struct{}

// NewTextResource creates a new text resource instance
func NewTextResource() *TextResource {
	return &TextResource{}
}

// Definition returns the resource definition for text elements
func (r *TextResource) Definition() providers.ResourceDefinition {
	return providers.ResourceDefinition{
		Type:        "text",
		Name:        "Text",
		Description: "Standalone text element for labels and annotations",
		Category:    "basic",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "Text content to display",
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
					"description": "Text width",
					"default":     100,
					"minimum":     10,
				},
				"height": map[string]interface{}{
					"type":        "number",
					"description": "Text height",
					"default":     30,
					"minimum":     10,
				},
				"fontSize": map[string]interface{}{
					"type":        "integer",
					"description": "Font size in points",
					"default":     12,
					"minimum":     6,
					"maximum":     72,
				},
				"fontFamily": map[string]interface{}{
					"type":        "string",
					"description": "Font family",
					"default":     "Arial",
				},
				"fontColor": map[string]interface{}{
					"type":        "string",
					"description": "Font color",
					"default":     "#000000",
				},
				"fontStyle": map[string]interface{}{
					"type":        "string",
					"description": "Font style",
					"default":     "normal",
					"enum":        []string{"normal", "bold", "italic", "bold italic"},
				},
				"textAlign": map[string]interface{}{
					"type":        "string",
					"description": "Text alignment",
					"default":     "center",
					"enum":        []string{"left", "center", "right"},
				},
				"verticalAlign": map[string]interface{}{
					"type":        "string",
					"description": "Vertical alignment",
					"default":     "middle",
					"enum":        []string{"top", "middle", "bottom"},
				},
				"fillColor": map[string]interface{}{
					"type":        "string",
					"description": "Background color (optional)",
					"default":     "",
				},
				"strokeColor": map[string]interface{}{
					"type":        "string",
					"description": "Border color (optional)",
					"default":     "",
				},
				"strokeWidth": map[string]interface{}{
					"type":        "number",
					"description": "Border width",
					"default":     0,
					"minimum":     0,
				},
			},
			"required": []string{"label"},
		},
		Examples: []providers.ResourceExample{
			{
				Name:        "Simple Text",
				Description: "Basic text label",
				Config: map[string]interface{}{
					"label":     "Sample Text",
					"x":         100,
					"y":         50,
					"fontSize":  14,
					"fontColor": "#333333",
				},
			},
			{
				Name:        "Styled Text Box",
				Description: "Text with background and border",
				Config: map[string]interface{}{
					"label":       "Important Note",
					"x":           200,
					"y":           100,
					"width":       150,
					"height":      40,
					"fontSize":    16,
					"fontStyle":   "bold",
					"textAlign":   "center",
					"fillColor":   "#FFF3CD",
					"strokeColor": "#856404",
					"strokeWidth": 2,
				},
			},
		},
	}
}

// Validate validates text parameters
func (r *TextResource) Validate(params map[string]interface{}) error {
	// Check required parameters
	if label, ok := params["label"]; !ok || label == "" {
		return &providers.ValidationError{
			Field:   "label",
			Message: "label text is required",
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

	// Validate text alignment if provided
	if textAlign, ok := params["textAlign"]; ok {
		validAligns := []string{"left", "center", "right"}
		if align, ok := textAlign.(string); ok {
			valid := false
			for _, validAlign := range validAligns {
				if align == validAlign {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "textAlign",
					Message: "invalid text alignment",
				}
			}
		}
	}

	// Validate vertical alignment if provided
	if verticalAlign, ok := params["verticalAlign"]; ok {
		validAligns := []string{"top", "middle", "bottom"}
		if align, ok := verticalAlign.(string); ok {
			valid := false
			for _, validAlign := range validAligns {
				if align == validAlign {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "verticalAlign",
					Message: "invalid vertical alignment",
				}
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
