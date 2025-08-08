package templates

import (
	"github.com/LederWorks/hippodamus/pkg/schema"
)

// TextTemplate handles template generation for text elements
type TextTemplate struct{}

// NewTextTemplate creates a new text template generator
func NewTextTemplate() *TextTemplate {
	return &TextTemplate{}
}

// Generate creates a schema.Element from text parameters
func (t *TextTemplate) Generate(params map[string]interface{}) (*schema.Element, error) {
	label := getStringParam(params, "label", "")
	x := getFloatParam(params, "x", 0)
	y := getFloatParam(params, "y", 0)
	width := getFloatParam(params, "width", 100)
	height := getFloatParam(params, "height", 30)

	// Style parameters
	fontSize := int(getFloatParam(params, "fontSize", 12))
	fontFamily := getStringParam(params, "fontFamily", "Arial")
	fontColor := getStringParam(params, "fontColor", "#000000")
	fontStyle := getStringParam(params, "fontStyle", "normal")
	textAlign := getStringParam(params, "textAlign", "center")
	verticalAlign := getStringParam(params, "verticalAlign", "middle")
	fillColor := getStringParam(params, "fillColor", "")
	strokeColor := getStringParam(params, "strokeColor", "")
	strokeWidth := getFloatParam(params, "strokeWidth", 0)

	return &schema.Element{
		Type: schema.ElementTypeText,
		Properties: schema.ElementProperties{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
			Label:  label,
		},
		Style: schema.Style{
			FontSize:      fontSize,
			FontFamily:    fontFamily,
			FontColor:     fontColor,
			FontStyle:     fontStyle,
			TextAlign:     textAlign,
			VerticalAlign: verticalAlign,
			FillColor:     fillColor,
			StrokeColor:   strokeColor,
			StrokeWidth:   strokeWidth,
		},
	}, nil
}
