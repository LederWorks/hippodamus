# Schema Validation Integration Script
# This script can be called from hippo.ps1 to validate schemas during the build process

[CmdletBinding()]
param(
    [Parameter(HelpMessage="Skip schema validation")]
    [switch]$SkipValidation,
    
    [Parameter(HelpMessage="Provider directory path")]
    [string]$ProvidersDir = "providers",
    
    [Parameter(HelpMessage="Fail build on validation errors")]
    [switch]$FailOnError
)

if ($SkipValidation) {
    Write-Host "⏭️  Schema validation skipped" -ForegroundColor Yellow
    return
}

$ValidationScript = Join-Path $ProvidersDir "schema-validation" "scripts" "validate-schema.ps1"

if (-not (Test-Path $ValidationScript)) {
    Write-Host "⚠️  Schema validation script not found at: $ValidationScript" -ForegroundColor Yellow
    Write-Host "   Run this from the root of the hippodamus repository, or use -SkipValidation" -ForegroundColor Yellow
    if ($FailOnError) {
        exit 1
    }
    return
}

Write-Host "🔍 Running schema validation..." -ForegroundColor Cyan

try {
    # Run validation in quiet mode, fail fast
    & $ValidationScript -Quiet -FailFast
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Schema validation passed" -ForegroundColor Green
    } else {
        Write-Host "❌ Schema validation failed" -ForegroundColor Red
        if ($FailOnError) {
            throw "Schema validation failed - see output above"
        }
    }
}
catch {
    Write-Host "❌ Error running schema validation: $($_.Exception.Message)" -ForegroundColor Red
    if ($FailOnError) {
        exit 1
    }
}
