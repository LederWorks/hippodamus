# Provider Organization Guidelines

## Recommended File Structure

### For Built-in Providers (within hippodamus repository)
```
providers/
├── <provider-name>/
│   ├── provider.go          # Main provider implementation & registry
│   ├── provider_test.go     # Provider-level tests
│   ├── resources/
│   │   ├── <resource1>.go       # Individual resource definition
│   │   ├── <resource1>_test.go  # Resource-specific tests
│   │   ├── <resource2>.go
│   │   ├── <resource2>_test.go
│   │   └── ...
│   ├── templates/
│   │   ├── <resource1>.go       # Template generation logic
│   │   ├── <resource1>_test.go  # Template tests
│   │   ├── <resource2>.go
│   │   └── ...
│   ├── validation/
│   │   ├── <resource1>.go       # Validation logic (if complex)
│   │   └── ...
│   ├── examples/
│   │   ├── <resource1>.yaml     # Usage examples
│   │   └── ...
│   ├── README.md            # Provider documentation
│   └── CHANGELOG.md         # Version history
```

### For External Providers (standalone repositories)
```
<provider-name>/
├── go.mod                   # Go module definition
├── go.sum
├── README.md               # Provider documentation
├── LICENSE                 # License file
├── CHANGELOG.md            # Version history
├── provider.go             # Main provider implementation
├── provider_test.go        # Provider-level tests
├── resources/
│   ├── <resource1>.go
│   ├── <resource1>_test.go
│   └── ...
├── templates/
│   ├── <resource1>.go
│   ├── <resource1>_test.go
│   └── ...
├── validation/             # Optional: for complex validation
│   └── ...
├── examples/
│   ├── simple/
│   │   └── diagram.yaml
│   ├── advanced/
│   │   └── complex-diagram.yaml
│   └── ...
├── docs/                   # Detailed documentation
│   ├── resources/
│   │   ├── <resource1>.md
│   │   └── ...
│   └── tutorials/
└── scripts/                # Build/test scripts
    ├── build.sh
    └── test.sh
```

## Resource Structure Standards

Each resource should follow this interface:

```go
// Resource definition
type ResourceDefinition struct {
    Type        string
    Name        string
    Description string
    Category    string
    Schema      map[string]interface{}
    Examples    []ResourceExample
}

// Resource implementation
type <ResourceName>Resource struct {
    provider *<Provider>Name
}

func New<ResourceName>Resource(provider *<Provider>Name) *<ResourceName>Resource {
    return &<ResourceName>Resource{provider: provider}
}

func (r *<ResourceName>Resource) Definition() ResourceDefinition { ... }
func (r *<ResourceName>Resource) Validate(params map[string]interface{}) error { ... }
func (r *<ResourceName>Resource) GenerateTemplate(params map[string]interface{}) (*schema.Element, error) { ... }
```

## Benefits of This Structure

1. **Modularity**: Each resource is self-contained
2. **Testability**: Individual resource testing
3. **Maintainability**: Easy to modify single resources
4. **Extensibility**: Simple to add new resources
5. **Documentation**: Clear resource-specific docs
6. **Portability**: Easy to move resources between providers
7. **Team Development**: Multiple developers can work on different resources

## Naming Conventions

- **Files**: snake_case (e.g., `shape_element.go`)
- **Types**: PascalCase (e.g., `ShapeElementResource`)
- **Functions**: PascalCase for public, camelCase for private
- **Constants**: UPPER_SNAKE_CASE
- **Resource Types**: lowercase with hyphens (e.g., `shape-element`)

## Versioning Strategy

- **Semantic Versioning**: MAJOR.MINOR.PATCH
- **API Compatibility**: Maintain backward compatibility within major versions
- **Schema Evolution**: Additive changes only in minor versions
- **Migration Guides**: Document breaking changes between major versions
