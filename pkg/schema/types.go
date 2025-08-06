package schema

import "time"

// DiagramConfig represents the root YAML configuration for a draw.io diagram
type DiagramConfig struct {
	Version   string        `yaml:"version" json:"version"`
	Metadata  Metadata      `yaml:"metadata" json:"metadata"`
	Templates []TemplateRef `yaml:"templates,omitempty" json:"templates,omitempty"`
	Diagram   Diagram       `yaml:"diagram" json:"diagram"`
}

// Metadata contains information about the diagram
type Metadata struct {
	Title       string    `yaml:"title" json:"title"`
	Description string    `yaml:"description,omitempty" json:"description,omitempty"`
	Author      string    `yaml:"author,omitempty" json:"author,omitempty"`
	Created     time.Time `yaml:"created,omitempty" json:"created,omitempty"`
	Modified    time.Time `yaml:"modified,omitempty" json:"modified,omitempty"`
	Tags        []string  `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// TemplateRef references an external template file
type TemplateRef struct {
	Name     string `yaml:"name" json:"name"`
	Template string `yaml:"template" json:"template"`
}

// Diagram represents the main diagram structure
type Diagram struct {
	Pages      []Page            `yaml:"pages" json:"pages"`
	Properties DiagramProperties `yaml:"properties,omitempty" json:"properties,omitempty"`
}

// DiagramProperties contains global diagram settings
type DiagramProperties struct {
	Grid       GridSettings       `yaml:"grid,omitempty" json:"grid,omitempty"`
	Background BackgroundSettings `yaml:"background,omitempty" json:"background,omitempty"`
	Scale      float64            `yaml:"scale,omitempty" json:"scale,omitempty"`
}

// GridSettings configures the diagram grid
type GridSettings struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Size    int    `yaml:"size,omitempty" json:"size,omitempty"`
	Color   string `yaml:"color,omitempty" json:"color,omitempty"`
}

// BackgroundSettings configures the diagram background
type BackgroundSettings struct {
	Color string `yaml:"color,omitempty" json:"color,omitempty"`
	Image string `yaml:"image,omitempty" json:"image,omitempty"`
}

// Page represents a single page in the diagram
type Page struct {
	ID         string         `yaml:"id" json:"id"`
	Name       string         `yaml:"name" json:"name"`
	Layers     []Layer        `yaml:"layers,omitempty" json:"layers,omitempty"`
	Elements   []Element      `yaml:"elements,omitempty" json:"elements,omitempty"`
	Properties PageProperties `yaml:"properties,omitempty" json:"properties,omitempty"`
}

// PageProperties contains page-specific settings
type PageProperties struct {
	Width      int    `yaml:"width,omitempty" json:"width,omitempty"`
	Height     int    `yaml:"height,omitempty" json:"height,omitempty"`
	Background string `yaml:"background,omitempty" json:"background,omitempty"`
}

// Layer represents a layer within a page
type Layer struct {
	ID       string    `yaml:"id" json:"id"`
	Name     string    `yaml:"name" json:"name"`
	Visible  bool      `yaml:"visible" json:"visible"`
	Locked   bool      `yaml:"locked,omitempty" json:"locked,omitempty"`
	Elements []Element `yaml:"elements,omitempty" json:"elements,omitempty"`
}

// Element represents any drawable element (shape, connector, text, etc.)
type Element struct {
	Type       ElementType       `yaml:"type" json:"type"`
	ID         string            `yaml:"id,omitempty" json:"id,omitempty"`
	Name       string            `yaml:"name,omitempty" json:"name,omitempty"`
	Template   string            `yaml:"template,omitempty" json:"template,omitempty"`
	Properties ElementProperties `yaml:"properties" json:"properties"`
	Style      Style             `yaml:"style,omitempty" json:"style,omitempty"`
	Children   []Element         `yaml:"children,omitempty" json:"children,omitempty"`
	Tags       []string          `yaml:"tags,omitempty" json:"tags,omitempty"`

	// Nesting configuration
	Nesting NestingConfig `yaml:"nesting,omitempty" json:"nesting,omitempty"`
}

// ElementType defines the type of element
type ElementType string

// Element type constants define the different types of elements that can be used
const (
	ElementTypeShape     ElementType = "shape"
	ElementTypeConnector ElementType = "connector"
	ElementTypeText      ElementType = "text"
	ElementTypeGroup     ElementType = "group"
	ElementTypeContainer ElementType = "container"
	ElementTypeSwimLane  ElementType = "swimlane"
	ElementTypeTemplate  ElementType = "template"
)

// NestingConfig defines how children should be nested within a parent element
type NestingConfig struct {
	Mode          NestingMode `yaml:"mode,omitempty" json:"mode,omitempty"`                   // How children are arranged
	AutoResize    bool        `yaml:"autoResize,omitempty" json:"autoResize,omitempty"`       // Auto-resize parent to fit children
	Padding       Padding     `yaml:"padding,omitempty" json:"padding,omitempty"`             // Padding around children
	Spacing       float64     `yaml:"spacing,omitempty" json:"spacing,omitempty"`             // Spacing between children
	Arrangement   Arrangement `yaml:"arrangement,omitempty" json:"arrangement,omitempty"`     // How children are arranged
	ChildDefaults *Element    `yaml:"childDefaults,omitempty" json:"childDefaults,omitempty"` // Default properties for children
}

// NestingMode defines how nesting behavior works
type NestingMode string

// Nesting mode constants
const (
	NestingModeContainer NestingMode = "container" // Default container behavior
	NestingModeGroup     NestingMode = "group"     // Group behavior (children move with parent)
	NestingModeSwimLane  NestingMode = "swimlane"  // Swimlane behavior
	NestingModeAutomatic NestingMode = "automatic" // Automatically determine based on element type
)

// Arrangement defines how children are arranged within a parent
type Arrangement string

// Arrangement constants
const (
	ArrangementVertical   Arrangement = "vertical"   // Stack children vertically
	ArrangementHorizontal Arrangement = "horizontal" // Stack children horizontally
	ArrangementGrid       Arrangement = "grid"       // Arrange in a grid
	ArrangementFree       Arrangement = "free"       // Free positioning
)

// Padding defines padding around nested content
type Padding struct {
	Top    float64 `yaml:"top,omitempty" json:"top,omitempty"`
	Right  float64 `yaml:"right,omitempty" json:"right,omitempty"`
	Bottom float64 `yaml:"bottom,omitempty" json:"bottom,omitempty"`
	Left   float64 `yaml:"left,omitempty" json:"left,omitempty"`
}

// ElementProperties contains position, size, and other element properties
type ElementProperties struct {
	// Position and size
	X      float64 `yaml:"x" json:"x"`
	Y      float64 `yaml:"y" json:"y"`
	Width  float64 `yaml:"width,omitempty" json:"width,omitempty"`
	Height float64 `yaml:"height,omitempty" json:"height,omitempty"`
	Z      int     `yaml:"z,omitempty" json:"z,omitempty"`

	// Content
	Label string `yaml:"label,omitempty" json:"label,omitempty"`
	Value string `yaml:"value,omitempty" json:"value,omitempty"`

	// Shape-specific
	Shape     string `yaml:"shape,omitempty" json:"shape,omitempty"`
	ShapeType string `yaml:"shapeType,omitempty" json:"shapeType,omitempty"`

	// Connector-specific
	Source     string     `yaml:"source,omitempty" json:"source,omitempty"`
	Target     string     `yaml:"target,omitempty" json:"target,omitempty"`
	SourcePort string     `yaml:"sourcePort,omitempty" json:"sourcePort,omitempty"`
	TargetPort string     `yaml:"targetPort,omitempty" json:"targetPort,omitempty"`
	Waypoints  []Waypoint `yaml:"waypoints,omitempty" json:"waypoints,omitempty"`

	// Group/Container-specific
	Collapsible bool `yaml:"collapsible,omitempty" json:"collapsible,omitempty"`
	Collapsed   bool `yaml:"collapsed,omitempty" json:"collapsed,omitempty"`

	// Custom properties for templates
	Custom map[string]interface{} `yaml:"custom,omitempty" json:"custom,omitempty"`
}

// Waypoint represents a point in a connector path
type Waypoint struct {
	X float64 `yaml:"x" json:"x"`
	Y float64 `yaml:"y" json:"y"`
}

// Style defines visual styling for elements
type Style struct {
	// Fill
	FillColor   string  `yaml:"fillColor,omitempty" json:"fillColor,omitempty"`
	FillOpacity float64 `yaml:"fillOpacity,omitempty" json:"fillOpacity,omitempty"`

	// Stroke
	StrokeColor     string  `yaml:"strokeColor,omitempty" json:"strokeColor,omitempty"`
	StrokeWidth     float64 `yaml:"strokeWidth,omitempty" json:"strokeWidth,omitempty"`
	StrokeOpacity   float64 `yaml:"strokeOpacity,omitempty" json:"strokeOpacity,omitempty"`
	StrokeDashArray string  `yaml:"strokeDashArray,omitempty" json:"strokeDashArray,omitempty"`

	// Text
	FontFamily            string `yaml:"fontFamily,omitempty" json:"fontFamily,omitempty"`
	FontSize              int    `yaml:"fontSize,omitempty" json:"fontSize,omitempty"`
	FontColor             string `yaml:"fontColor,omitempty" json:"fontColor,omitempty"`
	FontStyle             string `yaml:"fontStyle,omitempty" json:"fontStyle,omitempty"`                         // bold, italic, underline
	TextAlign             string `yaml:"textAlign,omitempty" json:"textAlign,omitempty"`                         // left, center, right
	VerticalAlign         string `yaml:"verticalAlign,omitempty" json:"verticalAlign,omitempty"`                 // top, middle, bottom
	LabelPosition         string `yaml:"labelPosition,omitempty" json:"labelPosition,omitempty"`                 // left, center, right
	VerticalLabelPosition string `yaml:"verticalLabelPosition,omitempty" json:"verticalLabelPosition,omitempty"` // top, middle, bottom

	// Shape-specific
	Rounded  bool    `yaml:"rounded,omitempty" json:"rounded,omitempty"`
	Shadow   bool    `yaml:"shadow,omitempty" json:"shadow,omitempty"`
	Glass    bool    `yaml:"glass,omitempty" json:"glass,omitempty"`
	Sketch   bool    `yaml:"sketch,omitempty" json:"sketch,omitempty"`
	Rotation float64 `yaml:"rotation,omitempty" json:"rotation,omitempty"`

	// Custom style properties
	Custom map[string]string `yaml:"custom,omitempty" json:"custom,omitempty"`
}

// IconConfig defines icon properties for elements
type IconConfig struct {
	Type        string  `yaml:"type,omitempty" json:"type,omitempty"`               // "shape" or "image"
	Shape       string  `yaml:"shape,omitempty" json:"shape,omitempty"`             // Shape name or image path
	FillColor   string  `yaml:"fillColor,omitempty" json:"fillColor,omitempty"`     // Fill color
	StrokeColor string  `yaml:"strokeColor,omitempty" json:"strokeColor,omitempty"` // Stroke color
	Position    string  `yaml:"position,omitempty" json:"position,omitempty"`       // Position relative to element
	Size        float64 `yaml:"size,omitempty" json:"size,omitempty"`               // Icon size
}

// GroupConfig defines the visual and behavioral properties of a template group
type GroupConfig struct {
	// Visual properties
	Properties ElementProperties `yaml:"properties,omitempty" json:"properties,omitempty"`
	Style      Style             `yaml:"style,omitempty" json:"style,omitempty"`

	// Container behavior
	AutoResize  bool        `yaml:"autoResize,omitempty" json:"autoResize,omitempty"`   // Auto-resize to fit children
	Padding     Padding     `yaml:"padding,omitempty" json:"padding,omitempty"`         // Padding around children
	Spacing     float64     `yaml:"spacing,omitempty" json:"spacing,omitempty"`         // Spacing between children
	Arrangement Arrangement `yaml:"arrangement,omitempty" json:"arrangement,omitempty"` // How children are arranged

	// Optional icon
	Icon *IconConfig `yaml:"icon,omitempty" json:"icon,omitempty"`

	// Optional child elements (like current templates)
	Children []Element `yaml:"children,omitempty" json:"children,omitempty"`
}

// Template represents a reusable template - now always a group that can contain other elements
type Template struct {
	Name         string       `yaml:"name" json:"name"`
	Description  string       `yaml:"description,omitempty" json:"description,omitempty"`
	Version      string       `yaml:"version,omitempty" json:"version,omitempty"`
	Dependencies []Dependency `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Parameters   []Parameter  `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Group        GroupConfig  `yaml:"group" json:"group"` // Every template defines a group
}

// Dependency defines a template dependency relationship
type Dependency struct {
	Name         string `yaml:"name" json:"name"`                                   // Logical name for the dependency
	Type         string `yaml:"type" json:"type"`                                   // Template type required
	Required     bool   `yaml:"required,omitempty" json:"required,omitempty"`       // Whether dependency is mandatory
	Description  string `yaml:"description,omitempty" json:"description,omitempty"` // Human-readable description
	Relationship string `yaml:"relationship" json:"relationship"`                   // "parent", "peer", "child", "ancestor"
	Multiple     bool   `yaml:"multiple,omitempty" json:"multiple,omitempty"`       // Allow multiple instances
}

// Parameter defines a template parameter
type Parameter struct {
	Name        string      `yaml:"name" json:"name"`
	Type        string      `yaml:"type" json:"type"` // string, number, boolean, color
	Default     interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Required    bool        `yaml:"required,omitempty" json:"required,omitempty"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
}
