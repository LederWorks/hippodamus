# Hippodamus Template Converter & Tester

This PowerShell script provides comprehensive template conversion and functional testing for the Hippodamus project.

## Features

### 🔄 Template Conversion
- **Automatic Detection**: Finds templates still using the old `elements:` format
- **Smart Conversion**: Handles both `container` and `shape` element types
- **Dependency Injection**: Automatically adds custom parent dependencies
- **Dry-run Mode**: Preview changes without modifying files

### 🧪 Functional Testing
- **Automatic Test Generation**: Creates test configurations for each template
- **End-to-end Validation**: Tests actual DrawIO file generation
- **XML Validation**: Verifies output files are valid XML
- **Structured Results**: Organizes test outputs in mirrored directory structure

## Directory Structure

The script expects and creates the following structure:
```
hippodamus/
├── templates/           # Template definitions
│   ├── aws/
│   ├── azure/
│   ├── gcp/
│   └── ...
├── config/             # Test configurations (auto-generated)
│   ├── aws/
│   ├── azure/
│   ├── gcp/
│   └── ...
├── results/            # Generated DrawIO files
│   ├── aws/
│   ├── azure/
│   ├── gcp/
│   └── ...
└── scripts/
    └── hippodamus-converter.ps1
```

## Usage

### Basic Conversion
```powershell
# Convert all remaining templates
.\scripts\hippodamus-converter.ps1

# Dry-run mode (preview only)
.\scripts\hippodamus-converter.ps1 -DryRun

# Skip testing after conversion
.\scripts\hippodamus-converter.ps1 -SkipTesting
```

### Testing Only
```powershell
# Test all templates without conversion
.\scripts\hippodamus-converter.ps1 -TestOnly

# Test with custom paths
.\scripts\hippodamus-converter.ps1 -TestOnly -ConfigDir "test-configs" -ResultsDir "test-results"
```

### Custom Paths
```powershell
# Use custom directories
.\scripts\hippodamus-converter.ps1 -TemplatesDir "my-templates" -ConfigDir "my-configs" -ResultsDir "my-results"

# Use custom hippodamus executable
.\scripts\hippodamus-converter.ps1 -HippodamusExecutable ".\build\hippo.exe"
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `TemplatesDir` | String | `"templates"` | Directory containing template files |
| `ConfigDir` | String | `"config"` | Directory for test configurations |
| `ResultsDir` | String | `"results"` | Directory for generated DrawIO files |
| `HippodamusExecutable` | String | `".\hippo.exe"` | Path to hippodamus executable |
| `DryRun` | Switch | `false` | Preview changes without modifying files |
| `SkipTesting` | Switch | `false` | Skip functional testing after conversion |
| `TestOnly` | Switch | `false` | Only run tests, skip conversion |

## What the Script Does

### 1. Template Conversion
- Scans for templates with `elements:` sections
- Converts to new `group:` format based on element type:
  - **Container elements**: Preserves properties, style, nesting, children
  - **Shape elements**: Converts to group with `autoResize: false` and `arrangement: "free"`
- Adds custom parent dependency for maximum flexibility
- Validates conversion success

### 2. Test Configuration Generation
- Automatically creates test configs for each template
- Parses template parameters and generates appropriate test values:
  - **String parameters**: Generated as "Test{ParameterName}"
  - **Color parameters**: Set to "#FF0000"
  - **Number parameters**: Set to "100"
  - **Parameters with defaults**: Uses default values
- Maintains directory structure matching templates

### 3. Functional Testing
- Runs hippodamus executable with each test configuration
- Validates DrawIO file generation
- Checks XML validity of output files
- Provides comprehensive test summary
- Organizes results in structured directories

## Example Output

```
🚀 Hippodamus Template Converter & Tester
===========================================

🔍 Finding templates that still have 'elements:' sections...
Template directory: templates

✅ All templates have been converted to the new group format!

🧪 Starting functional testing since all templates are converted...
Templates: templates
Configs: config
Results: results
Executable: .\hippo.exe

📝 Creating test config: config\aws\aws-account.yaml
🧪 Testing template: aws-account.yaml...
  ✅ Test passed - DrawIO file generated successfully
  ✅ Generated DrawIO file is valid XML

🎯 Testing Summary:
  Total templates tested: 34
  ✅ Passed: 34
  ❌ Failed: 0

🎉 All conversions validated successfully!
```

## Prerequisites

1. **Hippodamus executable**: Build the project first with `go build -o hippo.exe ./cmd/hippodamus`
2. **PowerShell**: Windows PowerShell 5.1+ or PowerShell Core 7+
3. **File permissions**: Write access to template, config, and results directories

## Error Handling

The script provides detailed error reporting:
- **Conversion errors**: Shows specific template and error details
- **Test failures**: Reports which templates failed testing and why
- **Missing dependencies**: Clear messages about missing executable or directories
- **XML validation**: Identifies invalid output files

## Integration

This script is designed to be used in CI/CD pipelines:
- Returns appropriate exit codes (0 for success, 1 for failure)
- Provides structured output suitable for automated processing
- Supports dry-run mode for validation workflows
- Can be run in test-only mode for regression testing
