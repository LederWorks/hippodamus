# Hippo - Build and Test Script

A streamlined PowerShell script that builds the Hippodamus application and tests all templates with intelligent defaults.

## Quick Start

```powershell
# Test with existing binary (default)
.\scripts\hippo.ps1

# Build only, skip tests
.\scripts\hippo.ps1 -SkipTest

# Build and test everything
.\scripts\hippo.ps1 -Build

# Test specific provider
.\scripts\hippo.ps1 -Provider azure

# Clean build with verbose output
.\scripts\hippo.ps1 -Build -Clean -Verbose
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Build` | Switch | `false` | Build the hippo.exe binary and run tests |
| `SkipTest` | Switch | `false` | Build only, skip testing (implies -Build) |
| `Provider` | String | `""` | Test specific provider only |
| `Binary` | String | `".\hippo.exe"` | Path to hippo binary |
| `ProvidersDir` | String | `"providers"` | Providers directory |
| `Clean` | Switch | `false` | Remove old binaries before build |
| `Verbose` | Switch | `false` | Show detailed output |

## Behavior Logic

The script has intelligent defaults based on the parameters provided:

- **No parameters**: `.\scripts\hippo.ps1` → Test only (uses existing binary)
- **SkipTest only**: `.\scripts\hippo.ps1 -SkipTest` → Build only (auto-enables -Build)
- **Build only**: `.\scripts\hippo.ps1 -Build` → Build and test
- **Build + SkipTest**: `.\scripts\hippo.ps1 -Build -SkipTest` → Build only

## What It Does

### Build Mode
- Validates Go installation
- Removes existing binaries (if `-Clean`)
- Compiles Go source into `hippo.exe`
- Tests the binary functionality
- Reports build size and status

### Test Mode (Default)
- Finds all providers with templates
- Auto-generates test configs for each template
- Runs hippo binary against each template
- Validates DrawIO XML output
- Reports pass/fail summary with counts

### Provider Testing
Test specific provider templates by using the `-Provider` parameter.

## Examples

```powershell
# Default: test all templates with existing binary
.\scripts\hippo.ps1

# Build only, no testing
.\scripts\hippo.ps1 -SkipTest

# Build and test everything
.\scripts\hippo.ps1 -Build

# Clean build only
.\scripts\hippo.ps1 -SkipTest -Clean

# Test AWS provider only
.\scripts\hippo.ps1 -Provider aws

# Build and test AWS only with verbose output
.\scripts\hippo.ps1 -Build -Provider aws -Verbose

# Custom binary location
.\scripts\hippo.ps1 -Binary ".\dist\hippo.exe"

# Build with clean and verbose, then test
.\scripts\hippo.ps1 -Build -Clean -Verbose
```

## Output Examples

### Build Only (`-SkipTest`)
```
Hippo - Build
======================
🔨 Building hippo...
ℹ️ Removed existing .\hippo.exe
✅ Built successfully: 5.3 MB
✅ Binary test passed: Hippodamus v1.0.0
✅ All operations completed successfully
```

### Test Only (Default)
```
Hippo - Test
======================
🧪 Testing templates...
ℹ️ Testing aws: 5 templates
ℹ️ Testing azure: 5 templates
ℹ️ Testing gcp: 5 templates
ℹ️ Testing kubernetes: 5 templates
ℹ️ Testing generic: 5 templates
ℹ️ Testing github: 2 templates
ℹ️ Testing azuredevops: 7 templates
✅ Tests: 34/34 passed
✅ All operations completed successfully
```

### Build and Test (`-Build`)
```
Hippo - Build and Test
======================
🔨 Building hippo...
✅ Built successfully: 5.3 MB
✅ Binary test passed: Hippodamus v1.0.0
🧪 Testing templates...
ℹ️ Testing aws: 5 templates
[... testing output ...]
✅ Tests: 34/34 passed
✅ All operations completed successfully
```

## Requirements

- **Go**: Latest version for building
- **PowerShell**: 5.1+ or PowerShell Core 7+
- **Directory Structure**: `providers/` with template subdirectories

## Error Handling

- Returns exit code 0 on success, 1 on failure
- Shows clear error messages with status emojis
- Stops on build failure before testing
- Continues testing all providers even if some fail

## File Generation

The script automatically creates:
- Test configurations in `providers/{provider}/configs/`
- DrawIO results in `providers/{provider}/results/`
- Directory structure as needed

Test configs are generated with:
- Basic YAML structure
- Template parameter detection
- Default test values for all parameters

## Status Emojis

- ✅ Success operations
- ❌ Error conditions
- ⚠️ Warning messages
- 🔨 Build operations
- 🧪 Test operations
- ℹ️ Informational messages

## Common Use Cases

1. **Development Workflow**: `.\scripts\hippo.ps1 -Build` - Build and validate everything
2. **CI/CD Pipeline**: `.\scripts\hippo.ps1 -SkipTest` followed by `.\scripts\hippo.ps1` - Separate build and test phases
3. **Provider Development**: `.\scripts\hippo.ps1 -Provider myProvider` - Test specific provider changes
4. **Release Preparation**: `.\scripts\hippo.ps1 -Build -Clean -Verbose` - Clean build with full details
