# Development Guide

## Project Structure

```
hippodamus/
├── cmd/
│   └── hippodamus/           # Main application entry point
│       └── main.go
├── pkg/
│   ├── schema/               # YAML schema definitions
│   │   └── types.go
│   ├── drawio/              # Draw.io XML generation
│   │   └── generator.go
│   └── templates/           # Template processing
│       └── processor.go
├── examples/                # Example configurations
│   ├── infrastructure.yaml  # Simple infrastructure diagram
│   ├── microservices.yaml  # Complex microservices architecture
│   ├── simple.yaml         # Basic diagram without templates
│   └── templates/          # Template definitions
│       ├── server.yaml
│       ├── database.yaml
│       ├── microservice.yaml
│       ├── container.yaml
│       └── loadbalancer.yaml
├── schemas/                 # JSON schemas for validation
│   └── diagram-config.schema.json
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Building and Testing

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)

### Building

```bash
# Build for current platform
go build -o hippodamus ./cmd/hippodamus

# Or using Make
make build

# Build for all platforms
make build-all
```

### Testing

```bash
# Run tests
go test -v ./...

# Or using Make
make test

# Run examples
make examples

# Validate examples
make validate
```

## Architecture

### Core Components

1. **Schema Package** (`pkg/schema/`)
   - Defines Go structs for YAML configuration
   - Handles type definitions for diagrams, pages, layers, elements
   - Supports templates, styling, and custom properties

2. **Draw.io Package** (`pkg/drawio/`)
   - Converts schema objects to draw.io XML format
   - Handles XML generation and serialization
   - Supports all major draw.io element types

3. **Templates Package** (`pkg/templates/`)
   - Processes template definitions and references
   - Applies template variables using Go's text/template
   - Merges template properties with element definitions

### Data Flow

1. **Load Configuration**: Parse YAML input file into schema structures
2. **Load Templates**: Load and parse referenced template files
3. **Process Templates**: Apply templates to elements, substitute variables
4. **Generate XML**: Convert processed schema to draw.io XML format
5. **Write Output**: Serialize XML to output file

## YAML Schema

### Basic Structure

```yaml
version: "1.0"
metadata:
  title: "Diagram Title"
  description: "Diagram description"
  author: "Author Name"
  tags: ["tag1", "tag2"]

templates:
  - name: "template-name"
    template: "path/to/template.yaml"

diagram:
  properties:
    grid:
      enabled: true
      size: 10
    background:
      color: "#ffffff"
  
  pages:
    - id: "page1"
      name: "Page Name"
      layers:
        - id: "layer1"
          name: "Layer Name"
          elements:
            - type: "shape"
              id: "element1"
              template: "template-name"
              properties:
                x: 100
                y: 100
                width: 120
                height: 80
                label: "Element Label"
              style:
                fillColor: "#color"
                strokeColor: "#color"
```

### Element Types

- **shape**: Basic shapes (rectangles, circles, etc.)
- **connector**: Lines and arrows connecting elements
- **text**: Text-only elements
- **group**: Groups of elements
- **container**: Container elements that can hold children
- **swimlane**: Swimlane elements for process diagrams

### Template System

Templates allow you to define reusable element configurations:

```yaml
name: "server"
description: "Server template"
parameters:
  - name: "serverType"
    type: "string"
    default: "generic"
  - name: "fillColor"
    type: "color"
    default: "#E3F2FD"

elements:
  - type: "shape"
    properties:
      width: 120
      height: 80
      label: "{{.label}}"
    style:
      fillColor: "{{.fillColor}}"
      strokeColor: "{{.strokeColor}}"
```

## Extending the System

### Adding New Element Types

1. Add the new type to `ElementType` enum in `pkg/schema/types.go`
2. Add handling in `generateElement` function in `pkg/drawio/generator.go`
3. Create a new generation function (e.g., `generateNewTypeCell`)

### Adding New Style Properties

1. Add properties to `Style` struct in `pkg/schema/types.go`
2. Update `generateElementStyle` function in `pkg/drawio/generator.go`
3. Update JSON schema in `schemas/diagram-config.schema.json`

### Adding Template Features

1. Extend `Template` or `Parameter` structs in `pkg/schema/types.go`
2. Update template processing in `pkg/templates/processor.go`
3. Add template variable handling in `applyTemplateVariables`

## Output Formats

Hippodamus supports two output formats:

1. **`.drawio` format** (default): The native draw.io format, identical to files created by diagrams.net
2. **`.xml` format**: Standard XML format with the same structure

Both formats contain identical XML content and are fully compatible with draw.io/diagrams.net.

## Draw.io XML Format

The generated XML follows the draw.io/diagrams.net format:

```xml
<mxfile>
  <diagram id="page1" name="Page Name">
    <mxGraphModel>
      <root>
        <mxCell id="0"/>
        <mxCell id="1" parent="0"/>
        <mxCell id="element1" value="Label" style="..." parent="1" vertex="1">
          <mxGeometry x="100" y="100" width="120" height="80" as="geometry"/>
        </mxCell>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
```

### Usage Examples

```bash
# Generate .drawio file (default)
./hippodamus -input diagram.yaml

# Generate .xml file
./hippodamus -input diagram.yaml -output diagram.xml

# Generate .drawio file explicitly
./hippodamus -input diagram.yaml -output diagram.drawio
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and validate examples
6. Submit a pull request

### Code Style

- Follow Go conventions
- Use `go fmt` for formatting
- Add comments for exported functions
- Write tests for new features
