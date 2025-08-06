package drawio

import (
	"encoding/xml"
	"fmt"
	"math"
	"strings"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

// DrawioDocument represents the root draw.io XML document
type DrawioDocument struct {
	XMLName  xml.Name        `xml:"mxfile"`
	Host     string          `xml:"host,attr"`
	Modified string          `xml:"modified,attr"`
	Agent    string          `xml:"agent,attr"`
	ETag     string          `xml:"etag,attr"`
	Version  string          `xml:"version,attr"`
	Type     string          `xml:"type,attr"`
	Diagram  []DrawioDiagram `xml:"diagram"`
}

// DrawioDiagram represents a single diagram/page
type DrawioDiagram struct {
	XMLName    xml.Name         `xml:"diagram"`
	ID         string           `xml:"id,attr"`
	Name       string           `xml:"name,attr"`
	GraphModel DrawioGraphModel `xml:"mxGraphModel"`
}

// DrawioGraphModel represents the graph model
type DrawioGraphModel struct {
	XMLName    xml.Name   `xml:"mxGraphModel"`
	DX         int        `xml:"dx,attr"`
	DY         int        `xml:"dy,attr"`
	Grid       int        `xml:"grid,attr"`
	GridSize   int        `xml:"gridSize,attr"`
	Guides     int        `xml:"guides,attr"`
	Tooltips   int        `xml:"tooltips,attr"`
	Connect    int        `xml:"connect,attr"`
	Arrows     int        `xml:"arrows,attr"`
	Fold       int        `xml:"fold,attr"`
	Page       int        `xml:"page,attr"`
	PageScale  float64    `xml:"pageScale,attr"`
	PageWidth  int        `xml:"pageWidth,attr"`
	PageHeight int        `xml:"pageHeight,attr"`
	Background string     `xml:"background,attr,omitempty"`
	Root       DrawioRoot `xml:"root"`
}

// DrawioRoot contains all cells
type DrawioRoot struct {
	XMLName xml.Name     `xml:"root"`
	Cells   []DrawioCell `xml:"mxCell"`
}

// DrawioCell represents a cell (shape, connector, etc.)
type DrawioCell struct {
	XMLName  xml.Name        `xml:"mxCell"`
	ID       string          `xml:"id,attr"`
	Value    string          `xml:"value,attr,omitempty"`
	Style    string          `xml:"style,attr,omitempty"`
	Parent   string          `xml:"parent,attr,omitempty"`
	Source   string          `xml:"source,attr,omitempty"`
	Target   string          `xml:"target,attr,omitempty"`
	Edge     string          `xml:"edge,attr,omitempty"`
	Vertex   string          `xml:"vertex,attr,omitempty"`
	Geometry *DrawioGeometry `xml:"mxGeometry,omitempty"`
}

// DrawioGeometry represents geometry information
type DrawioGeometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       float64  `xml:"x,attr,omitempty"`
	Y       float64  `xml:"y,attr,omitempty"`
	Width   float64  `xml:"width,attr,omitempty"`
	Height  float64  `xml:"height,attr,omitempty"`
	As      string   `xml:"as,attr"`
}

// Generator handles the conversion from schema to draw.io XML
type Generator struct {
	cellIDCounter int
}

// NewGenerator creates a new draw.io XML generator
func NewGenerator() *Generator {
	return &Generator{
		cellIDCounter: 0,
	}
}

// generateHierarchicalID creates a hierarchical ID path for an element
func (g *Generator) generateHierarchicalID(element *schema.Element, parentPath string) string {
	var identifier string

	// Use ID if provided, otherwise use Name
	if element.ID != "" {
		identifier = element.ID
	} else if element.Name != "" {
		identifier = element.Name
	} else {
		// This should not happen after validation, but provide fallback
		identifier = fmt.Sprintf("element-%d", g.cellIDCounter)
		g.cellIDCounter++
	}

	if parentPath == "" {
		return identifier
	}
	return parentPath + "/" + identifier
}

// validateElement validates that elements have required ID/Name fields
func (g *Generator) validateElement(element *schema.Element, isPageLevel bool) error {
	if isPageLevel {
		// Pages require both ID and Name
		if element.ID == "" || element.Name == "" {
			return fmt.Errorf("page elements must have both 'id' and 'name' fields")
		}
	} else {
		// Child elements require at least one of ID or Name
		if element.ID == "" && element.Name == "" {
			return fmt.Errorf("child elements must have either 'id' or 'name' field")
		}
	}

	// Recursively validate children
	for i := range element.Children {
		if err := g.validateElement(&element.Children[i], false); err != nil {
			return fmt.Errorf("invalid child element: %w", err)
		}
	}

	return nil
}

// getElementDisplayName returns the display name for an element (used for visible text)
func (g *Generator) getElementDisplayName(element *schema.Element) string {
	if element.Name != "" {
		return element.Name
	}
	return element.ID
}

// Generate converts a DiagramConfig to draw.io XML
func (g *Generator) Generate(config *schema.DiagramConfig) (*DrawioDocument, error) {
	doc := &DrawioDocument{
		Host:     "app.diagrams.net",
		Modified: "",
		Agent:    "Hippodamus",
		ETag:     "",
		Version:  "24.7.17",
		Type:     "device",
		Diagram:  make([]DrawioDiagram, 0, len(config.Diagram.Pages)),
	}

	for _, page := range config.Diagram.Pages {
		// Validate page has both ID and Name
		if page.ID == "" || page.Name == "" {
			return nil, fmt.Errorf("page must have both 'id' and 'name' fields")
		}

		// Validate all page elements
		for i := range page.Elements {
			if err := g.validateElement(&page.Elements[i], true); err != nil {
				return nil, fmt.Errorf("invalid element in page %s: %w", page.ID, err)
			}
		}

		diagram, err := g.generatePage(&page, &config.Diagram.Properties)
		if err != nil {
			return nil, fmt.Errorf("failed to generate page %s: %w", page.ID, err)
		}
		doc.Diagram = append(doc.Diagram, *diagram)
	}

	return doc, nil
}

// generatePage converts a Page to a DrawioDiagram
func (g *Generator) generatePage(page *schema.Page, diagramProps *schema.DiagramProperties) (*DrawioDiagram, error) {
	diagram := &DrawioDiagram{
		ID:   page.ID,
		Name: page.Name,
		GraphModel: DrawioGraphModel{
			DX:         0,
			DY:         0,
			Grid:       1,
			GridSize:   10,
			Guides:     1,
			Tooltips:   1,
			Connect:    1,
			Arrows:     1,
			Fold:       1,
			Page:       1,
			PageScale:  1,
			PageWidth:  827,
			PageHeight: 1169,
			Root: DrawioRoot{
				Cells: make([]DrawioCell, 0),
			},
		},
	}

	// Apply diagram properties
	if diagramProps != nil {
		if diagramProps.Grid.Enabled {
			diagram.GraphModel.Grid = 1
			if diagramProps.Grid.Size > 0 {
				diagram.GraphModel.GridSize = diagramProps.Grid.Size
			}
		} else {
			diagram.GraphModel.Grid = 0
		}

		if diagramProps.Scale > 0 {
			diagram.GraphModel.PageScale = diagramProps.Scale
		}

		if diagramProps.Background.Color != "" {
			diagram.GraphModel.Background = diagramProps.Background.Color
		}
	}

	// Apply page properties
	if page.Properties.Width > 0 {
		diagram.GraphModel.PageWidth = page.Properties.Width
	}
	if page.Properties.Height > 0 {
		diagram.GraphModel.PageHeight = page.Properties.Height
	}
	if page.Properties.Background != "" {
		diagram.GraphModel.Background = page.Properties.Background
	}

	// Add default root cells
	diagram.GraphModel.Root.Cells = append(diagram.GraphModel.Root.Cells,
		DrawioCell{ID: "0"},
		DrawioCell{ID: "1", Parent: "0"},
	)

	// Process layers
	for _, layer := range page.Layers {
		layerCell := DrawioCell{
			ID:     layer.ID,
			Value:  layer.Name,
			Parent: "1",
			Vertex: "1",
			Style:  g.generateLayerStyle(&layer),
		}
		diagram.GraphModel.Root.Cells = append(diagram.GraphModel.Root.Cells, layerCell)

		// Process elements in layer
		for _, element := range layer.Elements {
			cells, err := g.generateElement(&element, layer.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to generate element %s in layer %s: %w", element.ID, layer.ID, err)
			}
			diagram.GraphModel.Root.Cells = append(diagram.GraphModel.Root.Cells, cells...)
		}
	}

	// Process page-level elements (not in layers)
	for _, element := range page.Elements {
		cells, err := g.generateElementWithPath(&element, "1", page.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate element %s: %w", g.getElementDisplayName(&element), err)
		}
		diagram.GraphModel.Root.Cells = append(diagram.GraphModel.Root.Cells, cells...)
	}

	return diagram, nil
}

// generateElement converts an Element to DrawioCells
func (g *Generator) generateElement(element *schema.Element, parentID string) ([]DrawioCell, error) {
	var cells []DrawioCell

	// Apply automatic positioning if the element has children and nesting configuration
	if len(element.Children) > 0 {
		g.applyAutomaticNesting(element)
	}

	switch element.Type {
	case schema.ElementTypeShape, schema.ElementTypeContainer:
		cell := g.generateShapeCell(element, parentID)
		cells = append(cells, cell)

		// Process children
		for _, child := range element.Children {
			childCells, err := g.generateElement(&child, element.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to generate child element %s: %w", child.ID, err)
			}
			cells = append(cells, childCells...)
		}

	case schema.ElementTypeConnector:
		cell := g.generateConnectorCell(element, parentID)
		cells = append(cells, cell)

	case schema.ElementTypeText:
		cell := g.generateTextCell(element, parentID)
		cells = append(cells, cell)

	case schema.ElementTypeGroup:
		cell := g.generateGroupCell(element, parentID)
		cells = append(cells, cell)

		// Process group members
		for _, child := range element.Children {
			childCells, err := g.generateElement(&child, element.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to generate group member %s: %w", child.ID, err)
			}
			cells = append(cells, childCells...)
		}

	case schema.ElementTypeSwimLane:
		cell := g.generateSwimLaneCell(element, parentID)
		cells = append(cells, cell)

		// Process swimlane children
		for _, child := range element.Children {
			childCells, err := g.generateElement(&child, element.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to generate swimlane child %s: %w", child.ID, err)
			}
			cells = append(cells, childCells...)
		}

	default:
		return nil, fmt.Errorf("unsupported element type: %s", element.Type)
	}

	return cells, nil
}

// generateElementWithPath converts an Element to DrawioCells using hierarchical IDs
func (g *Generator) generateElementWithPath(element *schema.Element, parentID string, parentPath string) ([]DrawioCell, error) {
	var cells []DrawioCell

	// Generate hierarchical ID for this element
	elementPath := g.generateHierarchicalID(element, parentPath)

	// Temporarily update the element's ID to use the hierarchical path
	originalID := element.ID
	element.ID = elementPath

	// Apply automatic positioning if the element has children and nesting configuration
	if len(element.Children) > 0 {
		g.applyAutomaticNesting(element)
	}

	switch element.Type {
	case schema.ElementTypeShape, schema.ElementTypeContainer:
		cell := g.generateShapeCell(element, parentID)
		cells = append(cells, cell)

		// Process children
		for _, child := range element.Children {
			childCells, err := g.generateElementWithPath(&child, element.ID, elementPath)
			if err != nil {
				return nil, fmt.Errorf("failed to generate child element %s: %w", g.getElementDisplayName(&child), err)
			}
			cells = append(cells, childCells...)
		}

	case schema.ElementTypeConnector:
		cell := g.generateConnectorCell(element, parentID)
		cells = append(cells, cell)

	case schema.ElementTypeText:
		cell := g.generateTextCell(element, parentID)
		cells = append(cells, cell)

	case schema.ElementTypeGroup:
		cell := g.generateGroupCell(element, parentID)
		cells = append(cells, cell)

		// Process group members
		for _, child := range element.Children {
			childCells, err := g.generateElementWithPath(&child, element.ID, elementPath)
			if err != nil {
				return nil, fmt.Errorf("failed to generate group member %s: %w", g.getElementDisplayName(&child), err)
			}
			cells = append(cells, childCells...)
		}

	case schema.ElementTypeSwimLane:
		cell := g.generateSwimLaneCell(element, parentID)
		cells = append(cells, cell)

		// Process swimlane children
		for _, child := range element.Children {
			childCells, err := g.generateElementWithPath(&child, element.ID, elementPath)
			if err != nil {
				return nil, fmt.Errorf("failed to generate swimlane child %s: %w", g.getElementDisplayName(&child), err)
			}
			cells = append(cells, childCells...)
		}

	default:
		return nil, fmt.Errorf("unsupported element type: %s", element.Type)
	}

	// Restore original ID
	element.ID = originalID

	return cells, nil
}

// generateShapeCell creates a shape cell
func (g *Generator) generateShapeCell(element *schema.Element, parentID string) DrawioCell {
	// Use label if provided, otherwise use the display name (name or id)
	displayText := element.Properties.Label
	if displayText == "" {
		displayText = g.getElementDisplayName(element)
	}

	cell := DrawioCell{
		ID:     element.ID,
		Value:  displayText,
		Style:  g.generateElementStyle(element),
		Parent: parentID,
		Vertex: "1",
		Geometry: &DrawioGeometry{
			X:      element.Properties.X,
			Y:      element.Properties.Y,
			Width:  element.Properties.Width,
			Height: element.Properties.Height,
			As:     "geometry",
		},
	}

	return cell
}

// generateConnectorCell creates a connector cell
func (g *Generator) generateConnectorCell(element *schema.Element, parentID string) DrawioCell {
	cell := DrawioCell{
		ID:     element.ID,
		Value:  element.Properties.Label,
		Style:  g.generateElementStyle(element),
		Parent: parentID,
		Source: element.Properties.Source,
		Target: element.Properties.Target,
		Edge:   "1",
		Geometry: &DrawioGeometry{
			As: "geometry",
		},
	}

	return cell
}

// generateTextCell creates a text cell
func (g *Generator) generateTextCell(element *schema.Element, parentID string) DrawioCell {
	cell := DrawioCell{
		ID:     element.ID,
		Value:  element.Properties.Label,
		Style:  g.generateElementStyle(element),
		Parent: parentID,
		Vertex: "1",
		Geometry: &DrawioGeometry{
			X:      element.Properties.X,
			Y:      element.Properties.Y,
			Width:  element.Properties.Width,
			Height: element.Properties.Height,
			As:     "geometry",
		},
	}

	return cell
}

// generateGroupCell creates a group cell
func (g *Generator) generateGroupCell(element *schema.Element, parentID string) DrawioCell {
	// Use label if provided, otherwise use the display name (name or id)
	displayText := element.Properties.Label
	if displayText == "" {
		displayText = g.getElementDisplayName(element)
	}

	cell := DrawioCell{
		ID:     element.ID,
		Value:  displayText,
		Style:  g.generateElementStyle(element),
		Parent: parentID,
		Vertex: "1",
		Geometry: &DrawioGeometry{
			X:      element.Properties.X,
			Y:      element.Properties.Y,
			Width:  element.Properties.Width,
			Height: element.Properties.Height,
			As:     "geometry",
		},
	}

	return cell
}

// generateSwimLaneCell creates a swimlane cell
func (g *Generator) generateSwimLaneCell(element *schema.Element, parentID string) DrawioCell {
	style := g.generateElementStyle(element)
	if style == "" {
		style = "swimlane;fontStyle=0;childLayout=stackLayout;horizontal=1;startSize=30;horizontalStack=0;resizeParent=1;resizeParentMax=0;resizeLast=0;collapsible=1;marginBottom=0;"
	}

	// Use label if provided, otherwise use the display name (name or id)
	displayText := element.Properties.Label
	if displayText == "" {
		displayText = g.getElementDisplayName(element)
	}

	cell := DrawioCell{
		ID:     element.ID,
		Value:  displayText,
		Style:  style,
		Parent: parentID,
		Vertex: "1",
		Geometry: &DrawioGeometry{
			X:      element.Properties.X,
			Y:      element.Properties.Y,
			Width:  element.Properties.Width,
			Height: element.Properties.Height,
			As:     "geometry",
		},
	}

	return cell
}

// generateLayerStyle generates style for a layer
func (g *Generator) generateLayerStyle(layer *schema.Layer) string {
	var styles []string

	if !layer.Visible {
		styles = append(styles, "visible=0")
	}

	if layer.Locked {
		styles = append(styles, "locked=1")
	}

	return strings.Join(styles, ";")
}

// generateElementStyle generates style string for an element
func (g *Generator) generateElementStyle(element *schema.Element) string {
	var styles []string

	// Add shape type
	if element.Properties.Shape != "" {
		styles = append(styles, "shape="+element.Properties.Shape)
	}

	// Add style properties
	if element.Style.FillColor != "" {
		styles = append(styles, "fillColor="+element.Style.FillColor)
	}

	if element.Style.StrokeColor != "" {
		styles = append(styles, "strokeColor="+element.Style.StrokeColor)
	}

	if element.Style.StrokeWidth > 0 {
		styles = append(styles, fmt.Sprintf("strokeWidth=%.1f", element.Style.StrokeWidth))
	}

	if element.Style.StrokeDashArray != "" {
		styles = append(styles, "strokeDashArray="+element.Style.StrokeDashArray)
		styles = append(styles, "dashed=1")
	}

	if element.Style.FontFamily != "" {
		styles = append(styles, "fontFamily="+element.Style.FontFamily)
	}

	if element.Style.FontSize > 0 {
		styles = append(styles, fmt.Sprintf("fontSize=%d", element.Style.FontSize))
	}

	if element.Style.FontColor != "" {
		styles = append(styles, "fontColor="+element.Style.FontColor)
	}

	if element.Style.FontStyle != "" {
		styles = append(styles, "fontStyle="+element.Style.FontStyle)
	}

	if element.Style.TextAlign != "" {
		styles = append(styles, "align="+element.Style.TextAlign)
	}

	if element.Style.VerticalAlign != "" {
		styles = append(styles, "verticalAlign="+element.Style.VerticalAlign)
	}

	if element.Style.LabelPosition != "" {
		styles = append(styles, "labelPosition="+element.Style.LabelPosition)
	}

	if element.Style.VerticalLabelPosition != "" {
		styles = append(styles, "verticalLabelPosition="+element.Style.VerticalLabelPosition)
	}

	if element.Style.Rounded {
		styles = append(styles, "rounded=1")
	}

	if element.Style.Shadow {
		styles = append(styles, "shadow=1")
	}

	if element.Style.Glass {
		styles = append(styles, "glass=1")
	}

	if element.Style.Sketch {
		styles = append(styles, "sketch=1")
	}

	if element.Style.Rotation != 0 {
		styles = append(styles, fmt.Sprintf("rotation=%.1f", element.Style.Rotation))
	}

	// Add custom styles
	for key, value := range element.Style.Custom {
		styles = append(styles, fmt.Sprintf("%s=%s", key, value))
	}

	// Handle specific element types
	switch element.Type {
	case schema.ElementTypeConnector:
		if len(styles) == 0 || !containsStyle(styles, "edgeStyle") {
			styles = append(styles, "edgeStyle=orthogonalEdgeStyle")
		}
		if !containsStyle(styles, "html") {
			styles = append(styles, "html=1")
		}
		if !containsStyle(styles, "jettySize") {
			styles = append(styles, "jettySize=auto")
		}
		if !containsStyle(styles, "orthogonalLoop") {
			styles = append(styles, "orthogonalLoop=1")
		}

	case schema.ElementTypeText:
		if !containsStyle(styles, "text") {
			styles = append(styles, "text")
		}
		if !containsStyle(styles, "html") {
			styles = append(styles, "html=1")
		}
	}

	return strings.Join(styles, ";")
}

// containsStyle checks if a style property is already present
func containsStyle(styles []string, property string) bool {
	for _, style := range styles {
		if strings.HasPrefix(style, property+"=") || style == property {
			return true
		}
	}
	return false
}

// nextID generates a unique ID for cells
func (g *Generator) nextID() string {
	g.cellIDCounter++
	return fmt.Sprintf("cell-%d", g.cellIDCounter)
}

// applyAutomaticNesting applies automatic positioning to children based on nesting configuration
func (g *Generator) applyAutomaticNesting(element *schema.Element) {
	if len(element.Children) == 0 {
		return
	}

	nesting := element.Nesting

	// Set default nesting mode based on element type if not specified
	if nesting.Mode == "" {
		switch element.Type {
		case schema.ElementTypeContainer:
			nesting.Mode = schema.NestingModeContainer
		case schema.ElementTypeGroup:
			nesting.Mode = schema.NestingModeGroup
		case schema.ElementTypeSwimLane:
			nesting.Mode = schema.NestingModeSwimLane
		default:
			nesting.Mode = schema.NestingModeAutomatic
		}
	}

	// Set default arrangement if not specified
	if nesting.Arrangement == "" {
		if len(element.Children) <= 4 {
			nesting.Arrangement = schema.ArrangementHorizontal
		} else {
			nesting.Arrangement = schema.ArrangementGrid
		}
	}

	// Set default spacing if not specified
	if nesting.Spacing == 0 {
		nesting.Spacing = 20
	}

	// Set default padding if not specified
	if nesting.Padding.Top == 0 && nesting.Padding.Right == 0 &&
		nesting.Padding.Bottom == 0 && nesting.Padding.Left == 0 {
		nesting.Padding = schema.Padding{
			Top:    30,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}
	}

	// Calculate positions for children
	g.calculateChildPositions(element, &nesting)

	// Auto-resize parent if enabled
	if nesting.AutoResize {
		g.autoResizeParent(element, &nesting)
	}
}

// calculateChildPositions calculates automatic positions for child elements
func (g *Generator) calculateChildPositions(parent *schema.Element, nesting *schema.NestingConfig) {
	if len(parent.Children) == 0 {
		return
	}

	contentX := nesting.Padding.Left
	contentY := nesting.Padding.Top

	switch nesting.Arrangement {
	case schema.ArrangementVertical:
		currentY := contentY
		for i := range parent.Children {
			child := &parent.Children[i]
			if child.Properties.Width == 0 {
				child.Properties.Width = 140 // Default width
			}
			if child.Properties.Height == 0 {
				child.Properties.Height = 60 // Default height
			}

			child.Properties.X = contentX
			child.Properties.Y = currentY
			currentY += child.Properties.Height + nesting.Spacing
		}

	case schema.ArrangementHorizontal:
		currentX := contentX
		for i := range parent.Children {
			child := &parent.Children[i]
			if child.Properties.Width == 0 {
				child.Properties.Width = 140
			}
			if child.Properties.Height == 0 {
				child.Properties.Height = 60
			}

			child.Properties.X = currentX
			child.Properties.Y = contentY
			currentX += child.Properties.Width + nesting.Spacing
		}

	case schema.ArrangementGrid:
		// Calculate grid dimensions
		childCount := len(parent.Children)
		cols := int(math.Ceil(math.Sqrt(float64(childCount))))
		if cols > 4 {
			cols = 4 // Maximum 4 columns
		}

		currentX := contentX
		currentY := contentY
		col := 0

		for i := range parent.Children {
			child := &parent.Children[i]
			if child.Properties.Width == 0 {
				child.Properties.Width = 140
			}
			if child.Properties.Height == 0 {
				child.Properties.Height = 60
			}

			child.Properties.X = currentX
			child.Properties.Y = currentY

			col++
			if col >= cols {
				col = 0
				currentX = contentX
				currentY += child.Properties.Height + nesting.Spacing
			} else {
				currentX += child.Properties.Width + nesting.Spacing
			}
		}

	case schema.ArrangementFree:
		// Don't modify positions for free arrangement
		break
	}
}

// autoResizeParent automatically resizes the parent to fit all children
func (g *Generator) autoResizeParent(parent *schema.Element, nesting *schema.NestingConfig) {
	if len(parent.Children) == 0 {
		return
	}

	var maxX, maxY float64

	// Find the bottom-right bounds of all children
	for _, child := range parent.Children {
		childRight := child.Properties.X + child.Properties.Width
		childBottom := child.Properties.Y + child.Properties.Height

		if childRight > maxX {
			maxX = childRight
		}
		if childBottom > maxY {
			maxY = childBottom
		}
	}

	// Add padding to the calculated size
	newWidth := maxX + nesting.Padding.Right
	newHeight := maxY + nesting.Padding.Bottom

	// Only increase size, don't shrink
	if newWidth > parent.Properties.Width {
		parent.Properties.Width = newWidth
	}
	if newHeight > parent.Properties.Height {
		parent.Properties.Height = newHeight
	}
}
