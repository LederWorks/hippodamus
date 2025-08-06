# Hippodamus - YAML to Draw.io XML Converter

Hippodamus is a Go application that transforms YAML configuration files into draw.io XML syntax, allowing you to programmatically create diagrams from structured configuration files.

## Features

- Transform YAML configurations to draw.io XML format
- Support for nested diagrams with parents, layers, and tags
- Template system for reusable diagram components
- Customizable YAML schema for different diagram types
- Support for complex diagram structures including shapes, connectors, and styling

## Usage

```bash
# Build the application
go build -o hippodamus ./cmd/hippodamus

# Convert YAML to draw.io format (default .drawio extension)
./hippodamus -input diagram.yaml -output diagram.drawio

# Convert to XML format
./hippodamus -input diagram.yaml -output diagram.xml

# Use templates
./hippodamus -input diagram.yaml -templates templates/ -output diagram.drawio

# Auto-detect output format based on input filename
./hippodamus -input diagram.yaml
```

## Supported Output Formats

- `.drawio` - Draw.io native format (default)
- `.xml` - Standard XML format

Both formats contain the same XML structure and are fully compatible with draw.io/diagrams.net.

## YAML Schema

The application supports a comprehensive YAML schema for defining draw.io diagrams:

```yaml
version: "1.0"
metadata:
  title: "Sample Diagram"
  description: "A sample diagram configuration"
  
templates:
  - name: "server"
    template: "templates/server.yaml"
    
diagram:
  pages:
    - id: "page1"
      name: "Main Page"
      layers:
        - id: "layer1"
          name: "Infrastructure"
          elements:
            - type: "shape"
              id: "server1"
              template: "server"
              properties:
                x: 100
                y: 100
                width: 120
                height: 80
                label: "Web Server"
```

## Project Structure

- `cmd/hippodamus/` - Main application entry point
- `pkg/schema/` - YAML schema definitions and Go structs
- `pkg/drawio/` - Draw.io XML generation logic
- `pkg/templates/` - Template processing system
- `examples/` - Example YAML configurations and templates
- `schemas/` - JSON schemas for YAML validation
