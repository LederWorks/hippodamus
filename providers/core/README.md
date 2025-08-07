# Core Provider

The Core Provider supplies fundamental diagram elements for the Hippodamus system. This provider follows the modular organization pattern recommended for all Hippodamus providers.

## Organization Structure

```
providers/core/
├── provider.go              # Main provider implementation
├── provider_modular_test.go # Provider-level tests
├── resources/               # Resource definitions
│   ├── shape.go            # Shape resource definition
│   ├── shape_test.go       # Shape resource tests
│   ├── connector.go        # Connector resource definition
│   └── connector_test.go   # Connector resource tests
├── templates/              # Template generators
│   ├── shape.go           # Shape template generation
│   ├── shape_test.go      # Shape template tests
│   ├── connector.go       # Connector template generation
│   └── connector_test.go  # Connector template tests
└── README.md              # This file
```

## Supported Resources

### Shape (`core-shape`)
Basic shape element with customizable appearance and properties.

**Schema**: Rectangle, ellipse, triangle, diamond, hexagon, cloud, cylinder
**Parameters**: label, shape, width, height, x, y, fillColor, strokeColor, strokeWidth, fontSize, fontStyle, rounded, shadow

### Connector (`core-connector`)  
Connection line between elements with arrows and styling.

**Schema**: Source/target elements with port specification
**Parameters**: source, target, sourcePort, targetPort, label, strokeColor, strokeWidth, strokeStyle, arrow

### Text (TODO)
Standalone text element for labels and annotations.

### Group (TODO)
Container element that groups related elements together.

### Swimlane (TODO)
Horizontal or vertical lane for organizing process flows.

## Usage Examples

### Basic Shape
```yaml
- id: "web-server"
  name: "Web Server"
  type: "template"
  template: "core-shape"
  parameters:
    label: "Web Server"
    shape: "rectangle"
    fillColor: "#E3F2FD"
    strokeColor: "#1976D2"
    width: 120
    height: 80
```

### Basic Connector
```yaml
- id: "connection"
  name: "Database Connection"
  type: "template"
  template: "core-connector"
  parameters:
    source: "web-server"
    target: "database"
    label: "queries"
    strokeColor: "#424242"
    strokeWidth: 2
```

## Development Guidelines

### Adding New Resources

1. **Create Resource Definition**: Add `resources/<resource_name>.go`
2. **Create Resource Tests**: Add `resources/<resource_name>_test.go`
3. **Create Template Generator**: Add `templates/<resource_name>.go`
4. **Create Template Tests**: Add `templates/<resource_name>_test.go`
5. **Update Provider**: Add resource to `provider_modular.go`
6. **Update Tests**: Add integration tests

### Resource Structure Pattern

```go
type <ResourceName>Resource struct{}

func New<ResourceName>Resource() *<ResourceName>Resource {
    return &<ResourceName>Resource{}
}

func (r *<ResourceName>Resource) Definition() providers.ResourceDefinition {
    // Return resource schema and examples
}

func (r *<ResourceName>Resource) Validate(params map[string]interface{}) error {
    // Validate parameters according to schema
}
```

### Template Structure Pattern

```go
type <ResourceName>Template struct{}

func New<ResourceName>Template() *<ResourceName>Template {
    return &<ResourceName>Template{}
}

func (t *<ResourceName>Template) Generate(params map[string]interface{}) (*schema.Element, error) {
    // Generate schema.Element from parameters
}
```

## Testing

Run tests for specific components:

```bash
# Test all resources
go test ./providers/core/resources/...

# Test all templates  
go test ./providers/core/templates/...

# Test specific resource
go test ./providers/core/resources -run TestShapeResource

# Test with verbose output
go test -v ./providers/core/...
```

## Migration Status

- ✅ **Shape**: Resource and template implemented with full tests
- ✅ **Connector**: Resource and template implemented with full tests
- ⏳ **Text**: TODO - Migrate from monolithic structure
- ⏳ **Group**: TODO - Migrate from monolithic structure  
- ⏳ **Swimlane**: TODO - Migrate from monolithic structure

## Benefits of Modular Structure

1. **Separation of Concerns**: Resource definition, validation, and template generation are separate
2. **Testability**: Each component can be tested independently
3. **Maintainability**: Easy to modify individual resources without affecting others
4. **Extensibility**: Simple to add new resources following the established pattern
5. **Team Development**: Multiple developers can work on different resources simultaneously
6. **Reusability**: Resources can be easily moved between providers or reused
