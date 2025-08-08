package templates

import (
	"github.com/LederWorks/hippodamus/pkg/schema"
)

// ConnectorTemplate handles template generation for connector elements
type ConnectorTemplate struct{}

// NewConnectorTemplate creates a new connector template generator
func NewConnectorTemplate() *ConnectorTemplate {
	return &ConnectorTemplate{}
}

// Generate creates a schema.Element from connector parameters
func (t *ConnectorTemplate) Generate(params map[string]interface{}) (*schema.Element, error) {
	source := getStringParam(params, "source", "")
	target := getStringParam(params, "target", "")
	sourcePort := getStringParam(params, "sourcePort", "right")
	targetPort := getStringParam(params, "targetPort", "left")
	label := getStringParam(params, "label", "")
	strokeColor := getStringParam(params, "strokeColor", "#424242")
	strokeStyle := getStringParam(params, "strokeStyle", "solid")
	arrow := getStringParam(params, "arrow", "target")
	strokeWidth := getFloatParam(params, "strokeWidth", 2)

	return &schema.Element{
		Type: schema.ElementTypeConnector,
		Properties: schema.ElementProperties{
			Source:     source,
			Target:     target,
			SourcePort: sourcePort,
			TargetPort: targetPort,
			Label:      label,
			Custom: map[string]interface{}{
				"strokeStyle": strokeStyle,
				"arrow":       arrow,
			},
		},
		Style: schema.Style{
			StrokeColor: strokeColor,
			StrokeWidth: strokeWidth,
		},
	}, nil
}
