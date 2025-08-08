# Element Type Configurations

This directory contains standalone configuration files for each element type. These configs create complete, self-contained draw.io objects that can be used independently or as template inputs for larger diagrams.

## 🎯 Purpose

- **Object-Oriented Design**: Each element type has its own configuration
- **Reusability**: Configs can be used as building blocks for larger diagrams
- **No Duplication**: Single source of truth for each element type
- **Template Ready**: Can be used as inputs for template systems

## 📁 Available Element Configs

### 1. **Shape Element** (`shape.yaml`)
- **Type**: `shape`
- **Features**: Rectangle with rounded corners, comprehensive styling
- **Use Cases**: Basic shapes, UI elements, process boxes
- **Properties**: Position, size, colors, typography, effects
- **Customizable**: Shape type, colors, text, effects

### 2. **Connector Element** (`connector.yaml`)
- **Type**: `connector`
- **Features**: Line connecting two shapes with waypoints
- **Use Cases**: Flow lines, relationships, process flows
- **Properties**: Source/target, ports, waypoints, styling
- **Customizable**: Arrow types, line styles, colors, curves

### 3. **Text Element** (`text.yaml`)
- **Type**: `text`
- **Features**: Rich text with multiple formatting options
- **Use Cases**: Labels, descriptions, annotations, titles
- **Properties**: Content, typography, alignment, background
- **Customizable**: Fonts, colors, alignment, background

### 4. **Container Element** (`container.yaml`)
- **Type**: `container`
- **Features**: Layout container with child elements
- **Use Cases**: Panels, sections, grouped content
- **Properties**: Layout management, child arrangement, styling
- **Customizable**: Layout type, padding, child defaults

### 5. **Swimlane Element** (`swimlane.yaml`)
- **Type**: `swimlane`
- **Features**: Process swimlane with flow steps
- **Use Cases**: Process flows, workflows, pipelines
- **Properties**: Horizontal layout, process steps, connectors
- **Customizable**: Step count, colors, process names

### 6. **Group Element** (`group.yaml`)
- **Type**: `group`
- **Features**: Logical grouping with grid arrangement
- **Use Cases**: Related items, collections, categories
- **Properties**: Grid layout, collapsible, visual grouping
- **Customizable**: Grid size, colors, arrangement

## 🚀 Usage Examples

### Standalone Usage
```bash
# Use a shape element directly
hippo -Config providers/schema-validation/configs/element-types/shape.yaml

# Use a container element directly
hippo -Config providers/schema-validation/configs/element-types/container.yaml
```

### As Template Inputs
```yaml
# In a larger diagram configuration
templates:
  - name: "basic-shape"
    template: "schema-validation/configs/element-types/shape.yaml"
  - name: "process-flow"
    template: "schema-validation/configs/element-types/swimlane.yaml"

diagram:
  pages:
    - id: "main-page"
      name: "Main Diagram"
      elements:
        - type: "template"
          template: "basic-shape"
          properties:
            x: 100
            y: 100
            custom:
              label: "Custom Shape"
              fillColor: "#FF5722"
```

### As Building Blocks
```yaml
# Import multiple element types
diagram:
  pages:
    - id: "composite-page"
      name: "Composite Diagram"
      elements:
        # Use shape config as base
        - type: "shape"
          # ... properties from shape.yaml ...
          
        # Use container config as base  
        - type: "container"
          # ... properties from container.yaml ...
```

## 🔧 Customization

Each config is designed to be easily customizable:

### 1. **Properties Level**
- Change position (`x`, `y`)
- Modify size (`width`, `height`)
- Update content (`label`, `value`)

### 2. **Style Level**
- Colors (`fillColor`, `strokeColor`)
- Typography (`fontFamily`, `fontSize`)
- Effects (`shadow`, `rounded`)

### 3. **Custom Properties**
- Add template parameters
- Include metadata
- Extend functionality

## 📋 Configuration Structure

Each element config follows this pattern:

```yaml
version: "1.0"
metadata:
  title: "[Element] Element Configuration"
  description: "Standalone [element] element..."
  tags: ["[element]", "standalone", "reusable"]

diagram:
  properties:
    # Standard diagram properties
    
  pages:
    - id: "[element]-page"
      name: "[Element] Element"
      elements:
        - type: "[element]"
          id: "primary-[element]"
          properties:
            # Element-specific properties
          style:
            # Element-specific styling
          # Additional element features
```

## 🎨 Styling Consistency

All configs use consistent styling patterns:

- **Primary Colors**: Blue tones for shapes
- **Secondary Colors**: Complementary colors for variety
- **Typography**: Clear, readable fonts
- **Effects**: Subtle shadows and rounding
- **Spacing**: Consistent padding and margins

## 🔄 Template Integration

These configs are designed to work seamlessly with the template system:

1. **Parameter Support**: All configs accept custom parameters
2. **Override Friendly**: Properties can be easily overridden
3. **Composition Ready**: Multiple configs can be combined
4. **Inheritance**: Configs can extend each other

## 💡 Usage Examples

### Basic Usage
```yaml
# Reference an element configuration
elements:
  - id: "my_shape"
    type: "shape"
    extends: "@element-types/shape.yaml#basic_rectangle"
    position: { x: 100, y: 100 }
    content: "My Custom Shape"
```

### Composition Patterns
See the `examples/` directory for complete composition examples:

- `simple-composition.yaml` - Basic flowchart using element configs
- `object-oriented-composition.yaml` - Complex diagram with template inputs

### Template Integration
```yaml
# Import element configurations as template inputs
template_inputs:
  shapes: "@element-types/shape.yaml"
  connectors: "@element-types/connector.yaml"

# Use in diagram elements
elements:
  - id: "start"
    type: "shape"
    extends: "shapes.basic_rectangle"
    overrides:
      content: "Start Process"
      style:
        background_color: "#c8e6c9"
```

### Object-Oriented Composition
```yaml
version: "1.0"

# Import multiple element types
imports:
  shapes: "@element-types/shape.yaml"
  connectors: "@element-types/connector.yaml"
  containers: "@element-types/container.yaml"

diagram:
  elements:
    # Use container as foundation
    - id: "main_area"
      type: "container"
      extends: "containers.layout_container"
      children:
        # Add shapes as children
        - id: "process1"
          type: "shape"
          extends: "shapes.basic_rectangle"
          content: "Process 1"
        
        - id: "process2"
          type: "shape"
          extends: "shapes.basic_rectangle"
          content: "Process 2"
    
    # Connect with reusable connector
    - id: "flow"
      type: "connector"
      extends: "connectors.basic_connector"
      source: "process1"
      target: "process2"
```

## 📊 Validation

All element configs pass schema validation:

```bash
# Validate all element configs
cd providers/schema-validation
.\scripts\validate-schema.ps1 -Path "configs/element-types"
```

Use these element configs as your object-oriented building blocks for creating larger, more complex diagrams while maintaining consistency and avoiding duplication!
