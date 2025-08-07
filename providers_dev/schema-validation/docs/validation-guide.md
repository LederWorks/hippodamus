# Hippodamus Schema Validation Guide

This guide explains how to use the Schema Validation Provider to test and validate Hippodamus YAML configurations.

## 🎯 Purpose

The Schema Validation Provider helps ensure that:
- YAML configurations match the Go struct definitions in `types.go`
- Changes to the schema don't break existing functionality
- Developers understand all available schema features
- New configurations are valid before deployment

## 🚀 Getting Started

### Prerequisites
- PowerShell 5.1+ (Windows) or PowerShell Core 6+ (Cross-platform)
- Access to the Hippodamus repository

### Basic Usage

1. **Navigate to the provider directory:**
   ```powershell
   cd providers\schema-validation
   ```

2. **Run all validations:**
   ```powershell
   .\scripts\validate-schema.ps1
   ```

3. **Review the results:**
   - ✅ Green = Test passed
   - ❌ Red = Test failed
   - Summary shows total/passed/failed counts

## 📊 Test Types

### Positive Tests
Configurations that should pass validation:

```powershell
# Run only positive tests
.\scripts\validate-schema.ps1 -TestType positive
```

**Files tested:**
- `comprehensive-test.yaml` - Complete feature demonstration
- `minimal-test.yaml` - Minimal valid configuration

### Negative Tests  
Configurations that should fail validation:

```powershell
# Run only negative tests
.\scripts\validate-schema.ps1 -TestType negative
```

**Files tested:**
- `error-cases\invalid-nesting.yaml` - Invalid nesting modes
- `error-cases\missing-required.yaml` - Missing required fields  
- `error-cases\type-mismatch.yaml` - Wrong data types

### Edge Cases
Boundary conditions and unusual scenarios:

```powershell
# Run only edge case tests
.\scripts\validate-schema.ps1 -TestType edge-cases
```

**Files tested:**
- `edge-cases\boundary-conditions.yaml` - Extreme values, deep nesting

## 🔧 Command Line Options

### Basic Options
```powershell
# Test specific file
.\scripts\validate-schema.ps1 -Path "configs\comprehensive-test.yaml"

# Quiet mode (only show failures)
.\scripts\validate-schema.ps1 -Quiet

# Stop on first error
.\scripts\validate-schema.ps1 -FailFast

# Combine options
.\scripts\validate-schema.ps1 -TestType positive -Quiet
```

### Advanced Usage
```powershell
# Test all configs in a directory
.\scripts\validate-schema.ps1 -Path "configs\error-cases"

# Run from different directory
C:\> pwsh -File "C:\path\to\providers\schema-validation\scripts\validate-schema.ps1"
```

## 📝 Creating Test Cases

### Adding a Positive Test

1. **Create the YAML file:**
   ```yaml
   # configs/my-new-test.yaml
   version: "1.0"
   metadata:
     title: "My New Test"
   diagram:
     pages:
       - id: "test-page"
         name: "Test Page"
         elements:
           - type: "shape"
             id: "test-shape"
             properties:
               x: 100
               y: 100
               width: 100
               height: 80
               label: "Test Shape"
   ```

2. **Test it:**
   ```powershell
   .\scripts\validate-schema.ps1 -Path "configs\my-new-test.yaml"
   ```

3. **Verify it passes:**
   - Should show ✅ green checkmark
   - If it fails, fix the YAML structure

### Adding a Negative Test

1. **Create the YAML file:**
   ```yaml
   # configs/error-cases/my-error-test.yaml
   version: "1.0"
   diagram:
     pages:
       - id: "error-page"
         name: "Error Page"
         elements:
           - type: "invalid-type"  # This should fail
             id: "error-element"
             properties:
               x: 100
               y: 100
   ```

2. **Test it:**
   ```powershell
   .\scripts\validate-schema.ps1 -Path "configs\error-cases\my-error-test.yaml"
   ```

3. **Verify it fails correctly:**
   - Should show ❌ red X with error message
   - If it unexpectedly passes, review the validation logic

### Adding an Edge Case

1. **Create the YAML file:**
   ```yaml
   # configs/edge-cases/my-edge-test.yaml
   version: "1.0"
   metadata:
     title: "Edge Case: Zero-Sized Elements"
   diagram:
     pages:
       - id: "edge-page"
         name: "Edge Page"
         elements:
           - type: "shape"
             id: "zero-shape"
             properties:
               x: 100
               y: 100
               width: 0    # Edge case: zero width
               height: 0   # Edge case: zero height
               label: "Zero Size"
   ```

2. **Document the edge case:**
   Add a comment explaining what edge case is being tested

## 🔍 Validation Rules

The validation script checks these rules:

### Required Fields
- ✅ `version` field must be present
- ✅ `diagram` field must be present  
- ✅ Pages must have `id` and `name`
- ✅ Elements must have `type` and `properties`

### Valid Values
- ✅ Element types: shape, connector, text, group, container, swimlane, template
- ✅ Nesting modes: child, peer (old modes like container/group/swimlane are invalid)
- ✅ Arrangements: vertical, horizontal, grid, free

### Data Types
- ✅ Numeric fields must be numbers (not strings)
- ✅ Boolean fields must be true/false
- ✅ Color fields should be valid color strings

### Structure
- ✅ Valid YAML syntax
- ✅ Proper nesting structure
- ✅ No circular references

## 🐛 Troubleshooting

### Common Issues

**"Missing required 'version' field"**
```yaml
# ❌ Wrong - missing version
metadata:
  title: "My Diagram"

# ✅ Correct - includes version  
version: "1.0"
metadata:
  title: "My Diagram"
```

**"Invalid nesting mode found"**
```yaml
# ❌ Wrong - old nesting mode
nesting:
  mode: "container"

# ✅ Correct - new nesting mode
nesting:
  mode: "child"
```

**"Invalid element type"**
```yaml
# ❌ Wrong - invalid type
- type: "rectangle"

# ✅ Correct - valid type
- type: "shape"
  properties:
    shape: "rectangle"
```

### Getting Help

1. **Check the comprehensive test:**
   Look at `configs\comprehensive-test.yaml` for examples of correct usage

2. **Review the schema:**
   Check `pkg\schema\types.go` for the authoritative field definitions

3. **Run with verbose output:**
   Remove the `-Quiet` flag to see detailed validation messages

4. **Test incrementally:**
   Start with `minimal-test.yaml` and add features step by step

## 🔄 Integration with Development Workflow

### Before Making Schema Changes
```powershell
# Establish baseline
.\scripts\validate-schema.ps1 -Quiet
# Should pass all tests
```

### After Making Schema Changes
```powershell
# Check for regressions
.\scripts\validate-schema.ps1 -FailFast
# Fix any failures before committing
```

### Adding New Features
1. Update `types.go` with new schema fields
2. Add test cases demonstrating the new features
3. Update `comprehensive-test.yaml` if needed
4. Run full validation suite
5. Update documentation

### CI/CD Integration
```powershell
# In your build pipeline
if (Test-Path "providers\schema-validation\scripts\validate-schema.ps1") {
    & "providers\schema-validation\scripts\validate-schema.ps1" -Quiet -FailFast
    if ($LASTEXITCODE -ne 0) {
        throw "Schema validation failed - see output above"
    }
    Write-Host "✅ Schema validation passed" -ForegroundColor Green
}
```

This guide should help you effectively use the Schema Validation Provider to maintain high-quality, well-tested Hippodamus configurations.
