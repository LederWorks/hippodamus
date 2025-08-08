package templates

import (
	"github.com/LederWorks/hippodamus/pkg/schema"
)

// SwimlaneTemplate handles template generation for swimlane elements
type SwimlaneTemplate struct{}

// NewSwimlaneTemplate creates a new swimlane template generator
func NewSwimlaneTemplate() *SwimlaneTemplate {
	return &SwimlaneTemplate{}
}

// Generate creates a schema.Element from swimlane parameters
func (t *SwimlaneTemplate) Generate(params map[string]interface{}) (*schema.Element, error) {
	label := getStringParam(params, "label", "")
	x := getFloatParam(params, "x", 0)
	y := getFloatParam(params, "y", 0)
	width := getFloatParam(params, "width", 300)
	height := getFloatParam(params, "height", 200)

	// Swimlane specific properties
	orientation := getStringParam(params, "orientation", "horizontal")
	startSize := getFloatParam(params, "startSize", 30)
	childLayout := getStringParam(params, "childLayout", "stackLayout")
	collapsible := getBoolParam(params, "collapsible", true)
	collapsed := getBoolParam(params, "collapsed", false)

	// Style parameters
	fillColor := getStringParam(params, "fillColor", "#F8F9FA")
	strokeColor := getStringParam(params, "strokeColor", "#6C757D")
	strokeWidth := getFloatParam(params, "strokeWidth", 1)
	fontSize := int(getFloatParam(params, "fontSize", 12))
	fontStyle := getStringParam(params, "fontStyle", "bold")
	fontColor := getStringParam(params, "fontColor", "#000000")

	// Determine horizontal flag based on orientation
	horizontal := orientation == "horizontal"

	return &schema.Element{
		Type: schema.ElementTypeSwimLane,
		Properties: schema.ElementProperties{
			X:           x,
			Y:           y,
			Width:       width,
			Height:      height,
			Label:       label,
			Collapsible: collapsible,
			Collapsed:   collapsed,
			Custom: map[string]interface{}{
				"orientation": orientation,
				"startSize":   startSize,
				"childLayout": childLayout,
				"horizontal":  horizontal,
			},
		},
		Style: schema.Style{
			FillColor:   fillColor,
			StrokeColor: strokeColor,
			StrokeWidth: strokeWidth,
			FontSize:    fontSize,
			FontStyle:   fontStyle,
			FontColor:   fontColor,
		},
	}, nil
}
