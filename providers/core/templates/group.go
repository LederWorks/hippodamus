package templates

import (
	"github.com/LederWorks/hippodamus/pkg/schema"
)

// GroupTemplate handles template generation for group elements
type GroupTemplate struct{}

// NewGroupTemplate creates a new group template generator
func NewGroupTemplate() *GroupTemplate {
	return &GroupTemplate{}
}

// Generate creates a schema.Element from group parameters
func (t *GroupTemplate) Generate(params map[string]interface{}) (*schema.Element, error) {
	label := getStringParam(params, "label", "")
	x := getFloatParam(params, "x", 0)
	y := getFloatParam(params, "y", 0)
	width := getFloatParam(params, "width", 200)
	height := getFloatParam(params, "height", 150)

	// Container properties
	collapsible := getBoolParam(params, "collapsible", false)
	collapsed := getBoolParam(params, "collapsed", false)

	// Style parameters
	fillColor := getStringParam(params, "fillColor", "#F5F5F5")
	strokeColor := getStringParam(params, "strokeColor", "#CCCCCC")
	strokeWidth := getFloatParam(params, "strokeWidth", 1)
	strokeStyle := getStringParam(params, "strokeStyle", "solid")
	rounded := getBoolParam(params, "rounded", false)
	fontSize := int(getFloatParam(params, "fontSize", 12))
	fontStyle := getStringParam(params, "fontStyle", "bold")
	fontColor := getStringParam(params, "fontColor", "#000000")

	return &schema.Element{
		Type: schema.ElementTypeGroup,
		Properties: schema.ElementProperties{
			X:           x,
			Y:           y,
			Width:       width,
			Height:      height,
			Label:       label,
			Collapsible: collapsible,
			Collapsed:   collapsed,
		},
		Style: schema.Style{
			FillColor:   fillColor,
			StrokeColor: strokeColor,
			StrokeWidth: strokeWidth,
			Rounded:     rounded,
			FontSize:    fontSize,
			FontStyle:   fontStyle,
			FontColor:   fontColor,
			Custom: map[string]string{
				"strokeStyle": strokeStyle,
			},
		},
	}, nil
}
