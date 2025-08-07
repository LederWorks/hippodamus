# Schema Validation Provider

The Schema Validation Provider is a comprehensive testing and validation framework for the Hippodamus YAML schema defined in `pkg/schema/types.go`. This provider serves multiple purposes:

- **Schema Testing**: Validate that YAML configurations conform to the Go struct definitions
- **Regression Testing**: Ensure schema changes don't break existing functionality  
- **Documentation**: Provide working examples of all schema features
- **Developer Onboarding**: Help new developers understand the schema capabilities

## 🏗️ Provider Structure

```
schema-validation/
├── configs/                     # Test configurations
│   ├── comprehensive-test.yaml  # Complete feature demonstration
│   ├── minimal-test.yaml       # Minimal valid configuration
│   ├── error-cases/            # Invalid configurations (should fail)
│   │   ├── invalid-nesting.yaml
│   │   ├── missing-required.yaml
│   │   └── type-mismatch.yaml
│   └── edge-cases/             # Edge cases and boundary conditions
│       └── boundary-conditions.yaml
├── templates/                  # Validation template files
│   ├── validation-suite.yaml   # Comprehensive validation template
│   └── test-container.yaml     # Simple test container template
├── scripts/                    # Validation and testing scripts
│   └── validate-schema.ps1     # PowerShell validation script
├── docs/                      # Documentation
│   ├── schema-reference.md     # Auto-generated schema reference
│   └── validation-guide.md     # How to use this provider
├── results/                   # Test results and reports
└── README.md                  # This file
```

## 🚀 Quick Start

### Run All Validations

```powershell
# From the schema-validation directory
.\scripts\validate-schema.ps1
```

### Run Specific Test Types

```powershell
# Only positive tests (should pass)
.\scripts\validate-schema.ps1 -TestType positive

# Only negative tests (should fail)  
.\scripts\validate-schema.ps1 -TestType negative

# Only edge cases
.\scripts\validate-schema.ps1 -TestType edge-cases

# Quiet mode (only show failures)
.\scripts\validate-schema.ps1 -Quiet

# Fail fast (stop on first error)
.\scripts\validate-schema.ps1 -FailFast
```

### Test Individual Files

```powershell
# Test a specific configuration
.\scripts\validate-schema.ps1 -Path "configs\comprehensive-test.yaml"
```

## 📋 Test Categories

### 1. Positive Tests
These configurations should pass validation and demonstrate correct usage:

- **`comprehensive-test.yaml`**: Exercises all schema features including:
  - All element types (shape, connector, text, group, container, swimlane, template)
  - All nesting modes (child, peer)  
  - All arrangement types (vertical, horizontal, grid, free)
  - Complex styling and properties
  - Template usage and parameters
  - Multi-page layouts with layers

- **`minimal-test.yaml`**: Minimal valid configuration showing the absolute minimum required fields

### 2. Negative Tests (Error Cases)
These configurations should fail validation:

- **`invalid-nesting.yaml`**: Uses deprecated nesting modes (container, group, swimlane, automatic)
- **`missing-required.yaml`**: Missing required fields like version, page names, element IDs
- **`type-mismatch.yaml`**: Incorrect data types (strings for numbers, etc.)

### 3. Edge Cases
These test boundary conditions and unusual but valid scenarios:

- **`boundary-conditions.yaml`**: Tests extreme values, zero sizes, negative coordinates, deep nesting

## 🔧 Schema Features Tested

### Core Structure
- ✅ DiagramConfig with version, metadata, templates, diagram
- ✅ Metadata with all optional fields
- ✅ Template references
- ✅ Multi-page diagrams
- ✅ Layer organization

### Element Types
- ✅ **Shape**: Basic shapes with all properties and styling
- ✅ **Connector**: Lines between elements with waypoints
- ✅ **Text**: Text elements with typography settings
- ✅ **Group**: Grouping containers with child elements
- ✅ **Container**: Advanced containers with layout management
- ✅ **Swimlane**: Process swimlanes with flow elements
- ✅ **Template**: Template instances with parameters

### Nesting System
- ✅ **Child Mode**: Elements contained within parent (sub-grouping)
- ✅ **Peer Mode**: Elements that move together as peers (same-level grouping)
- ✅ **Auto-resize**: Automatic parent resizing to fit children
- ✅ **Padding**: Configurable padding around nested content
- ✅ **Spacing**: Configurable spacing between child elements

### Layout Arrangements  
- ✅ **Vertical**: Stack children vertically
- ✅ **Horizontal**: Stack children horizontally  
- ✅ **Grid**: Arrange children in a grid pattern
- ✅ **Free**: Absolute positioning of children

### Styling System
- ✅ **Fill**: Colors, opacity, gradients
- ✅ **Stroke**: Colors, width, opacity, dash patterns
- ✅ **Typography**: Font family, size, color, style, alignment
- ✅ **Shape Effects**: Rounded corners, shadows, glass effects
- ✅ **Custom Properties**: Extensible custom styling

### Template System
- ✅ **Parameters**: Typed parameters with defaults and validation
- ✅ **Dependencies**: Template dependency relationships
- ✅ **Group Configuration**: Template-based group definitions
- ✅ **Parameter Substitution**: Template parameter interpolation

## 🔍 Validation Rules

The validation script checks for:

1. **Required Fields**: Ensures all mandatory fields are present
2. **Valid Types**: Validates element types against allowed values
3. **Nesting Modes**: Ensures only `child` and `peer` modes are used
4. **Data Types**: Validates field types match schema expectations
5. **YAML Structure**: Ensures valid YAML syntax and structure

## 🧪 Integration with CI/CD

This provider can be integrated into your build pipeline:

```powershell
# In your build script (hippo.ps1)
& "providers\schema-validation\scripts\validate-schema.ps1" -Quiet -FailFast
if ($LASTEXITCODE -ne 0) {
    throw "Schema validation failed"
}
```

## 📝 Adding New Tests

### Adding Positive Tests
1. Create a new `.yaml` file in `configs/`
2. Ensure it follows the schema correctly
3. Run validation to confirm it passes

### Adding Negative Tests
1. Create a new `.yaml` file in `configs/error-cases/`
2. Include intentional schema violations
3. Run validation to confirm it fails as expected

### Adding Edge Cases
1. Create a new `.yaml` file in `configs/edge-cases/`
2. Test boundary conditions or unusual scenarios
3. Document what specific edge case is being tested

## 🔄 Schema Evolution

When modifying `pkg/schema/types.go`:

1. **Before Changes**: Run all validations to establish baseline
2. **After Changes**: Run validations to check for regressions
3. **Update Tests**: Add new test cases for new features
4. **Update Documentation**: Keep schema reference up to date

## 📊 Validation Reports

Validation results are displayed with clear indicators:

- ✅ **Green**: Test passed
- ❌ **Red**: Test failed  
- 📋 **Yellow**: Test category headers
- 🎉 **Celebration**: All tests passed

The script provides detailed error messages for failures and a comprehensive summary.

## 🤝 Contributing

When contributing to this provider:

1. **Test Coverage**: Ensure new schema features have corresponding tests
2. **Documentation**: Update this README for significant changes
3. **Validation**: Run the full validation suite before committing
4. **Examples**: Provide clear examples of new features

This provider ensures the Hippodamus schema remains robust, well-tested, and documented as it evolves.
