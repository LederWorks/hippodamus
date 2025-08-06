# Schema Validation Script for Hippodamus
# This script validates YAML configs against the types.go schema

[CmdletBinding()]
param(
    [Parameter(HelpMessage="Path to config file or directory to validate")]
    [string]$Path = ".",
    
    [Parameter(HelpMessage="Skip successful validation output")]
    [switch]$Quiet,
    
    [Parameter(HelpMessage="Stop on first error")]
    [switch]$FailFast,
    
    [Parameter(HelpMessage="Validate only specific test type")]
    [ValidateSet("all", "positive", "negative", "edge-cases")]
    [string]$TestType = "all"
)

Write-Host "🔍 Hippodamus Schema Validation" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan

$ValidationResults = @()
$TotalTests = 0
$PassedTests = 0
$FailedTests = 0

function Test-YamlStructure {
    param(
        [string]$FilePath,
        [bool]$ShouldPass = $true
    )
    
    $script:TotalTests++
    $TestName = Split-Path $FilePath -Leaf
    
    try {
        # Basic YAML parsing test
        $content = Get-Content $FilePath -Raw
        if (-not $content) {
            throw "File is empty"
        }
        
        # Check for required fields
        if ($content -notmatch 'version:') {
            throw "Missing required 'version' field"
        }
        
        if ($content -notmatch 'diagram:') {
            throw "Missing required 'diagram' field"
        }
        
        # Check for valid nesting modes
        if ($content -match 'mode:\s*"?(container|group|swimlane|automatic)"?') {
            throw "Invalid nesting mode found. Only 'child' and 'peer' are valid."
        }
        
        # Check for invalid element types (more specific matching to avoid false positives)
        $ValidTypes = @("shape", "connector", "text", "group", "container", "swimlane", "template")
        # Only match 'type:' that appears as a direct property of list items (elements)
        $typeMatches = [regex]::Matches($content, '(?:^|\n)\s*-[^\r\n]*\n(?:\s+[^\r\n]*\n)*\s+type:\s*"?([^""\s\r\n]+)"?', 'Multiline')
        foreach ($match in $typeMatches) {
            $type = $match.Groups[1].Value.Trim()
            if ($type -notin $ValidTypes) {
                throw "Invalid element type: '$type'. Valid types are: $($ValidTypes -join ', ')"
            }
        }
        
        # Check for missing required fields in pages
        if ($content -match 'pages:') {
            # Check if any page is missing name
            if ($content -match '- id:\s*[^#\r\n]+' -and $content -notmatch 'name:') {
                throw "Page missing required 'name' field"
            }
        }
        
        # Check for type mismatches - basic validation
        if ($content -match 'width:\s*"[^"]*"' -and $content -match 'width:\s*"[^\d\.]') {
            throw "Width field should be numeric, not string"
        }
        if ($content -match 'height:\s*"[^"]*"' -and $content -match 'height:\s*"[^\d\.]') {
            throw "Height field should be numeric, not string"
        }
        
        if ($ShouldPass) {
            $script:PassedTests++
            if (-not $Quiet) {
                Write-Host "✅ $TestName" -ForegroundColor Green
            }
            return @{ Status = "PASS"; File = $TestName; Error = $null }
        } else {
            $script:FailedTests++
            Write-Host "❌ $TestName - Expected to fail but passed" -ForegroundColor Red
            return @{ Status = "UNEXPECTED_PASS"; File = $TestName; Error = "Expected failure but validation passed" }
        }
    }
    catch {
        if ($ShouldPass) {
            $script:FailedTests++
            Write-Host "❌ $TestName - $($_.Exception.Message)" -ForegroundColor Red
            return @{ Status = "FAIL"; File = $TestName; Error = $_.Exception.Message }
        } else {
            $script:PassedTests++
            if (-not $Quiet) {
                Write-Host "✅ $TestName (Expected failure)" -ForegroundColor Green
            }
            return @{ Status = "EXPECTED_FAIL"; File = $TestName; Error = $_.Exception.Message }
        }
    }
}

# Determine what to validate
$ConfigsPath = Join-Path $PSScriptRoot ".." "configs"

# If specific path provided, validate that instead
if ($Path -ne ".") {
    if (Test-Path $Path) {
        if ((Get-Item $Path).PSIsContainer) {
            # Directory - validate all YAML files in it
            Write-Host "`n📋 Validating directory: $Path..." -ForegroundColor Yellow
            Get-ChildItem $Path -Filter "*.yaml" | ForEach-Object {
                $ValidationResults += Test-YamlStructure $_.FullName $true
            }
        } else {
            # Single file
            Write-Host "`n📋 Validating file: $Path..." -ForegroundColor Yellow
            $ValidationResults += Test-YamlStructure $Path $true
        }
    } else {
        Write-Host "❌ Path not found: $Path" -ForegroundColor Red
        exit 1
    }
} else {
    # Run standard test categories

    if ($TestType -eq "all" -or $TestType -eq "positive") {
        Write-Host "`n📋 Running Positive Tests (should pass)..." -ForegroundColor Yellow
    
    # Test comprehensive config
    $comprehensive = Join-Path $ConfigsPath "comprehensive-test.yaml"
    if (Test-Path $comprehensive) {
        $ValidationResults += Test-YamlStructure $comprehensive $true
    }
    
    # Test minimal config
    $minimal = Join-Path $ConfigsPath "minimal-test.yaml"
    if (Test-Path $minimal) {
        $ValidationResults += Test-YamlStructure $minimal $true
    }
}

if ($TestType -eq "all" -or $TestType -eq "negative") {
    Write-Host "`n📋 Running Negative Tests (should fail)..." -ForegroundColor Yellow
    
    # Test error cases
    $errorCases = Join-Path $ConfigsPath "error-cases"
    if (Test-Path $errorCases) {
        Get-ChildItem $errorCases -Filter "*.yaml" | ForEach-Object {
            $ValidationResults += Test-YamlStructure $_.FullName $false
            if ($FailFast -and $ValidationResults[-1].Status -eq "FAIL") {
                break
            }
        }
    }
}

if ($TestType -eq "all" -or $TestType -eq "edge-cases") {
    Write-Host "`n📋 Running Edge Case Tests..." -ForegroundColor Yellow
    
    # Test edge cases
    $edgeCases = Join-Path $ConfigsPath "edge-cases"
    if (Test-Path $edgeCases) {
        Get-ChildItem $edgeCases -Filter "*.yaml" | ForEach-Object {
            $ValidationResults += Test-YamlStructure $_.FullName $true
            if ($FailFast -and $ValidationResults[-1].Status -eq "FAIL") {
                break
            }
        }
    }
}
}

# Summary
Write-Host "`n📊 Validation Summary" -ForegroundColor Cyan
Write-Host "=====================" -ForegroundColor Cyan
Write-Host "Total Tests: $TotalTests" -ForegroundColor White
Write-Host "Passed: $PassedTests" -ForegroundColor Green
Write-Host "Failed: $FailedTests" -ForegroundColor $(if ($FailedTests -eq 0) { "Green" } else { "Red" })

if ($FailedTests -gt 0) {
    Write-Host "`n❌ Failed Tests:" -ForegroundColor Red
    $ValidationResults | Where-Object { $_.Status -eq "FAIL" -or $_.Status -eq "UNEXPECTED_PASS" } | ForEach-Object {
        Write-Host "  - $($_.File): $($_.Error)" -ForegroundColor Red
    }
    exit 1
} else {
    Write-Host "`n🎉 All validations passed!" -ForegroundColor Green
    exit 0
}
