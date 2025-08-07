# 🎨 Authentic Draw.io Element Type Configurations

This directory contains individual configurations for each **authentic** draw.io element type. These configurations are based on actual draw.io XML export structure and represent the real element types used in draw.io diagrams.

## 🎯 **Purpose**

- **Authentic Types**: Only element types that exist in actual draw.io XML
- **Object-Oriented Design**: Each element type has its own configuration
- **Reusability**: Configs can be used as building blocks for larger diagrams
- **No Duplication**: Single source of truth for each element type
- **Template Ready**: Can be used as inputs for template systems

## 📋 **Authentic Element Types**

Based on draw.io XML analysis, these are the **real** element types:

### 1. **Shape Element** (`shape.yaml`)
- **Type**: `shape` 
- **Draw.io Properties**: `vertex=1, connectable=1`
- **Features**: All geometric shapes, icons, container-like behavior
- **Use Cases**: Rectangles, circles, AWS/Azure icons, visual containers, tables
- **Shape Variants**: `rectangle`, `ellipse`, `triangle`, `diamond`, `hexagon`, `process`, `table`, `tableRow`, `partialRectangle`, `mxgraph.AWS4.RESOURCE`, `image`
- **Container Behavior**: Achieved with `nesting.mode: "child"` + `rounded: 1`

### 2. **Connector Element** (`connector.yaml`)
- **Type**: `connector`
- **Draw.io Properties**: `edge=1, vertex=0, connectable=0`
- **Features**: Source/target connections, waypoints, edge styling
- **Use Cases**: Arrows, lines, relationships, flows
- **Edge Styles**: `orthogonalEdgeStyle`, `none` (straight), curved
- **Required**: Must have `source` and `target` properties

### 3. **Text Element** (`text.yaml`)
- **Type**: `text`
- **Draw.io Properties**: `vertex=1, connectable=0`
- **Features**: Rich text with formatting options
- **Use Cases**: Labels, descriptions, annotations, standalone text
- **Properties**: Typography, alignment, wrapping
- **Styling**: Usually `strokeColor=none, fillColor=none`

### 4. **Group Element** (`group.yaml`)
- **Type**: `group`
- **Draw.io Properties**: `vertex=1, connectable=1`
- **Features**: Logical grouping with free child positioning
- **Use Cases**: Collections, logical groupings, free layouts
- **Nesting**: Always `mode: "child"` with `arrangement: "free"`
- **Appearance**: Usually `rounded: 0` (sharp edges)

### 5. **Swimlane Element** (`swimlane.yaml`)
- **Type**: `swimlane`
- **Draw.io Properties**: `vertex=1, connectable=1, horizontal=1, childLayout=stackLayout`
- **Features**: Process lanes with header and horizontal flow
- **Use Cases**: Process flows, workflows, pipelines
- **Layout**: Horizontal child arrangement with header space
- **Special**: `swimlaneHead=0` (no header section)

### 6. **Template Element** (referenced in configurations)
- **Type**: `template`
- **Draw.io Properties**: `vertex=1, connectable=1`
- **Features**: Template instance references with parameters
- **Use Cases**: Reusable components, parameterized elements
- **Properties**: Template references, parameter passing

## 🚫 **Removed Artificial Types**

**Important**: The following types were **removed** as they don't exist in actual draw.io XML:

- ~~`container`~~ - Use `shape` with `nesting.mode: "child"` instead

## 🔍 **Key Insight: Container Behavior**

**Container behavior** is achieved through configuration, not element type:

```yaml
# Container-like element (actually a shape)
type: "shape"
properties:
  shape: "rectangle"
style:
  rounded: 1              # Visual container appearance
  whiteSpace: "wrap"
  html: 1
nesting:
  mode: "child"           # Container behavior
  arrangement: "vertical" # Layout management
  spacing: 15
  padding: { top: 20, right: 15, bottom: 15, left: 15 }
children:
  - type: "shape"         # Child elements
    # ... child configuration
```

This matches actual draw.io XML where "containers" are shapes with specific styling and child arrangement properties.

## 🎨 **Nesting Modes**

All authentic element types can use any nesting mode:

| Nesting Mode | Purpose | When to Use |
|--------------|---------|-------------|
| `individual` | Standalone elements | Leaf nodes, simple content |
| `peer` | Elements that move together | Related items, grouped movement |
| `child` | Container behavior | Layout management, hierarchies |

## 🔧 **Shape Variants**

The `shape` element type supports many variants through the `shape` property:

### **Geometric Shapes**
```yaml
shape: "rectangle"    # Basic rectangle
shape: "ellipse"      # Circle/oval
shape: "triangle"     # Triangle
shape: "diamond"      # Diamond
shape: "hexagon"      # Hexagon
```

### **Process Shapes**
```yaml
shape: "process"      # Process box
```

### **Table Elements**
```yaml
shape: "table"              # Main table
shape: "tableRow"           # Row within table
shape: "partialRectangle"   # Cell within row
```

### **Cloud/Icon Shapes**
```yaml
shape: "mxgraph.AWS4.RESOURCE"  # AWS icons
shape: "image"                  # Azure/other icons
```

## 💡 **Usage Examples**

### **Basic Shape (Individual)**
```yaml
- type: "shape"
  properties:
    shape: "rectangle"
  nesting: { mode: "individual" }
```

### **Container-like Shape (Child)**
```yaml
- type: "shape"
  properties:
    shape: "rectangle"
  style:
    rounded: 1
    whiteSpace: "wrap"
    html: 1
  nesting:
    mode: "child"
    arrangement: "vertical"
  children:
    - type: "shape"
      # ... child configuration
```

### **Process Flow (Swimlane)**
```yaml
- type: "swimlane"
  style:
    horizontal: 1
    childLayout: "stackLayout"
  nesting:
    mode: "child"
    arrangement: "horizontal"
  children:
    - type: "shape"
      properties:
        shape: "process"
```

## 📊 **Validation**

All element configs pass schema validation:

```bash
# Validate all element configs
.\scripts\validate-schema.ps1 -Path "configs/element-types"

# Validate specific config
.\scripts\validate-schema.ps1 -Path "configs/element-types/shape.yaml"
```

## 🚀 **Template Integration**

These configs are designed to work seamlessly with the template system:

```yaml
# Import element configurations as template inputs
template_inputs:
  shapes: "@element-types/shape.yaml"
  connectors: "@element-types/connector.yaml"

# Use in diagram elements
elements:
  - type: "shape"
    extends: "shapes.primary-shape"
    overrides:
      content: "Custom Content"
      style:
        fillColor: "#custom-color"
```

Use these authentic element configs as your object-oriented building blocks for creating any draw.io diagram structure with proper nesting behavior that matches actual draw.io XML output!
