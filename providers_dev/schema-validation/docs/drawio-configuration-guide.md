# 🎯 Draw.io Element Configuration Guide

## Your Requirements Analysis

Based on your draw.io export analysis and the diagram you shared, here's the complete configuration system:

### 📋 **The Three Nesting Modes System**

| Nesting Mode | Purpose | When to Use | Draw.io Properties |
|--------------|---------|-------------|-------------------|
| **Individual** | Standalone elements | Leaf nodes, simple content | `vertex=1, connectable=1` |
| **Peer** | Elements that move together | Related items, grouped movement | `vertex=1, connectable=1` |
| **Child** | Container elements | Layout management, hierarchies | `vertex=1, connectable=1` |

---

## 🎨 **Style Configurations**

### **Shapes** (`vertex=1, connectable=1`)
```yaml
# Basic rectangle
shape: "rectangle"
rounded: 0                # Sharp corners
vertex: 1
parent: 1
connectable: 1

# Table structure
shape: "table"            # Main table
shape: "tableRow"         # Row within table  
shape: "partialRectangle" # Cell within row

# Cloud resources
shape: "mxgraph.AWS4.RESOURCE"  # AWS icons
shape: "image"                  # Azure icons with sketch=1
```

### **Text** (`vertex=1, connectable=0`)
```yaml
type: "text"
vertex: 1
parent: 1
connectable: 0            # Usually not connectable
whiteSpace: "wrap"
html: 1
strokeColor: "none"
fillColor: "none"
```

### **Swimlanes** (`vertex=1, connectable=1`)
```yaml
type: "swimlane"
vertex: 1
parent: 1
connectable: 1
swimlaneHead: 0           # No header section
horizontal: 1             # Horizontal layout
childLayout: "stackLayout"
```

### **Images/Icons** (`vertex=1, connectable=1`)
```yaml
shape: "image"
vertex: 1
parent: 1
connectable: 1
aspect: "fixed"           # Maintain aspect ratio
sketch: 1                 # Azure sketch style
```

### **Groups** (`vertex=1, connectable=1`)
```yaml
type: "group"
vertex: 1
parent: 1
connectable: 1
rounded: 0                # Usually sharp
```

### **Containers** (`vertex=1, connectable=1`)
```yaml
type: "container"
vertex: 1
parent: 1
connectable: 1
rounded: 0                # Sharp container
rounded: 1                # Rounded container
whiteSpace: "wrap"
html: 1
```

---

## 🔗 **Connector Properties**

### **Edge Styles** (`edge=1, vertex=0`)
```yaml
type: "connector"
edge: 1                   # Is an edge, not vertex
vertex: 0
parent: 1
source: "element-id"      # Required: source element
target: "element-id"      # Required: target element
connectable: 0            # Connectors not connectable

# Style variants:
edgeStyle: "orthogonalEdgeStyle"  # 90-degree angles
edgeStyle: "none"                 # Straight lines
```

---

## 🔧 **Special Properties Explained**

### **parent** (always required)
- **Purpose**: References the parent container/page
- **Usage**: `parent: 1` (page level) or `parent: "container-id"`
- **Required**: Every element must have a parent reference

### **vertex** (geometry vs connections)
- **vertex=1**: Element has geometry (x, y, width, height)
  - Shapes, containers, groups, swimlanes, text
- **vertex=0**: Element is a connection
  - Only connectors use vertex=0

### **connectable** (can be connected to)
- **connectable=1**: Other elements can connect to this
  - Shapes, containers, groups, swimlanes
- **connectable=0**: Cannot be connected to
  - Text elements, connectors, table cells

### **source/target** (connector endpoints)
- **source**: ID of element where connector starts
- **target**: ID of element where connector ends
- **Required**: Only for connectors

### **edge** (connection indicator)
- **edge=1**: Element is a connection line
- **Only used**: For connectors
- **Mutually exclusive**: with vertex=1

---

## 📐 **Practical Usage Patterns**

### **Pattern 1: AWS Infrastructure Table**
```yaml
# Container (child mode)
type: "container"
nesting: { mode: "child" }
children:
  # Table (individual mode)
  - type: "shape"
    shape: "table"
    nesting: { mode: "individual" }
    children:
      # Table rows
      - type: "shape"
        shape: "tableRow"
        children:
          # AWS icons in cells
          - type: "shape"
            shape: "mxgraph.AWS4.RESOURCE"
            nesting: { mode: "individual" }
```

### **Pattern 2: Microservices Peer Group**
```yaml
# Services that move together
- type: "shape"
  nesting: { mode: "peer", spacing: 20 }
  # Peer group 1

- type: "shape" 
  nesting: { mode: "peer", spacing: 20 }
  # Peer group 2 (same spacing)
```

### **Pattern 3: Process Swimlane**
```yaml
# Swimlane container (child mode)
type: "swimlane"
childLayout: "stackLayout"
nesting: { mode: "child", arrangement: "horizontal" }
children:
  # Process steps
  - type: "shape"
    shape: "process"
    nesting: { mode: "individual" }
```

---

## ✅ **Your Complete Configuration**

Based on your requirements, here's what you need:

1. **✅ Three nesting modes**: `individual`, `peer`, `child`
2. **✅ All draw.io shapes**: rectangle, table, AWS resources, Azure images
3. **✅ Container types**: container, swimlane, group
4. **✅ Connector support**: orthogonal, straight, curved
5. **✅ Proper properties**: vertex, parent, connectable, source, target, edge

The configurations in `drawio-element-mapping.yaml` and `practical-nesting-demo.yaml` give you everything you need to build any draw.io diagram structure with proper nesting behavior!

## 🚀 **Next Steps**

1. Use `practical-nesting-demo.yaml` as your template
2. Pick elements from `drawio-element-mapping.yaml` 
3. Combine using the three nesting modes
4. Test with: `.\scripts\validate-schema.ps1`
