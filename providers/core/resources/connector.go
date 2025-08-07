package resources

import (
	"github.com/LederWorks/hippodamus/pkg/providers"
)

// ConnectorResource defines the connector element resource
type ConnectorResource struct{}

// NewConnectorResource creates a new connector resource instance
func NewConnectorResource() *ConnectorResource {
	return &ConnectorResource{}
}

// Definition returns the resource definition for connector elements
func (r *ConnectorResource) Definition() providers.ResourceDefinition {
	return providers.ResourceDefinition{
		Type:        "connector",
		Name:        "Connector",
		Description: "Connection line between elements with arrows and styling",
		Category:    "basic",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"source": map[string]interface{}{
					"type":        "string",
					"description": "Source element ID",
				},
				"target": map[string]interface{}{
					"type":        "string",
					"description": "Target element ID",
				},
				"sourcePort": map[string]interface{}{
					"type":        "string",
					"description": "Source connection point",
					"default":     "right",
					"enum":        []string{"top", "right", "bottom", "left", "center"},
				},
				"targetPort": map[string]interface{}{
					"type":        "string",
					"description": "Target connection point",
					"default":     "left",
					"enum":        []string{"top", "right", "bottom", "left", "center"},
				},
				"label": map[string]interface{}{
					"type":        "string",
					"description": "Label for the connector",
					"default":     "",
				},
				"strokeColor": map[string]interface{}{
					"type":        "string",
					"description": "Line color",
					"default":     "#424242",
				},
				"strokeWidth": map[string]interface{}{
					"type":        "number",
					"description": "Line width",
					"default":     2,
					"minimum":     1,
				},
				"strokeStyle": map[string]interface{}{
					"type":        "string",
					"description": "Line style",
					"default":     "solid",
					"enum":        []string{"solid", "dashed", "dotted"},
				},
				"arrow": map[string]interface{}{
					"type":        "string",
					"description": "Arrow style",
					"default":     "target",
					"enum":        []string{"none", "source", "target", "both"},
				},
			},
			"required": []string{"source", "target"},
		},
		Examples: []providers.ResourceExample{
			{
				Name:        "Basic Connection",
				Description: "Simple connector between two elements",
				Config: map[string]interface{}{
					"source":      "element1",
					"target":      "element2",
					"sourcePort":  "right",
					"targetPort":  "left",
					"strokeColor": "#424242",
				},
			},
			{
				Name:        "Labeled Dashed Line",
				Description: "Dashed connector with label",
				Config: map[string]interface{}{
					"source":      "start",
					"target":      "end",
					"label":       "Data Flow",
					"strokeStyle": "dashed",
					"arrow":       "target",
				},
			},
		},
	}
}

// Validate validates connector parameters
func (r *ConnectorResource) Validate(params map[string]interface{}) error {
	// Check required parameters
	if source, ok := params["source"]; !ok || source == "" {
		return &providers.ValidationError{
			Field:   "source",
			Message: "source element ID is required",
		}
	}

	if target, ok := params["target"]; !ok || target == "" {
		return &providers.ValidationError{
			Field:   "target",
			Message: "target element ID is required",
		}
	}

	// Validate port values if provided
	validPorts := []string{"top", "right", "bottom", "left", "center"}

	if sourcePort, ok := params["sourcePort"]; ok {
		if port, ok := sourcePort.(string); ok {
			valid := false
			for _, validPort := range validPorts {
				if port == validPort {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "sourcePort",
					Message: "invalid source port value",
				}
			}
		}
	}

	if targetPort, ok := params["targetPort"]; ok {
		if port, ok := targetPort.(string); ok {
			valid := false
			for _, validPort := range validPorts {
				if port == validPort {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "targetPort",
					Message: "invalid target port value",
				}
			}
		}
	}

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

	// Validate arrow style if provided
	if arrow, ok := params["arrow"]; ok {
		validArrows := []string{"none", "source", "target", "both"}
		if arrowStr, ok := arrow.(string); ok {
			valid := false
			for _, validArrow := range validArrows {
				if arrowStr == validArrow {
					valid = true
					break
				}
			}
			if !valid {
				return &providers.ValidationError{
					Field:   "arrow",
					Message: "invalid arrow style",
				}
			}
		}
	}

	// Validate stroke width if provided
	if strokeWidth, ok := params["strokeWidth"]; ok {
		if w, ok := strokeWidth.(float64); ok && w < 1 {
			return &providers.ValidationError{
				Field:   "strokeWidth",
				Message: "stroke width must be at least 1",
			}
		}
	}

	return nil
}
