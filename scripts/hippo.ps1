# Hippo - Build and Test Script
param(
    [switch]$Build = $false,
    [switch]$SkipTest = $false,
    [string]$Provider = "",
    [string]$Binary = ".\hippo.exe",
    [string]$ProvidersDir = "providers",
    [switch]$Clean = $false,
    [switch]$Verbose = $false
)

# If no specific action, default to Test only
if (-not $Build -and -not $SkipTest) {
    $Test = $true
} elseif ($SkipTest -and -not $Build) {
    # If SkipTest is specified without Build, default to Build
    $Build = $true
    $Test = $false
} else {
    $Test = -not $SkipTest
}

function Write-Status {
    param([string]$Message, [string]$Status = "INFO")
    $emoji = switch ($Status) {
        "SUCCESS" { "✅" }
        "ERROR" { "❌" }
        "WARNING" { "⚠️" }
        "BUILD" { "🔨" }
        "TEST" { "🧪" }
        default { "ℹ️" }
    }
    Write-Host "$emoji $Message"
}

function Build-Hippo {
    param([string]$OutputPath, [bool]$Clean, [bool]$Verbose)
    
    Write-Status "Building hippo..." "BUILD"
    
    if ($Clean) {
        $binaries = @("hippo.exe", "hippodamus.exe", "main.exe")
        foreach ($binary in $binaries) {
            if (Test-Path $binary) {
                Remove-Item -Path $binary -Force
                Write-Status "Removed $binary"
            }
        }
    }
    
    if (Test-Path $OutputPath) {
        Remove-Item -Path $OutputPath -Force
        Write-Status "Removed existing $OutputPath"
    }
    
    # Check Go
    try {
        $goVersion = go version 2>&1
        if ($LASTEXITCODE -ne 0) { throw "Go not found" }
        if ($Verbose) { Write-Status "Go: $goVersion" }
    }
    catch {
        Write-Status "Go is not installed or not in PATH" "ERROR"
        return $false
    }
    
    # Build
    try {
        if ($Verbose) {
            $result = go build -v -o $OutputPath ./cmd/hippodamus 2>&1
        } else {
            $result = go build -o $OutputPath ./cmd/hippodamus 2>&1
        }
        
        if ($LASTEXITCODE -eq 0 -and (Test-Path $OutputPath)) {
            $fileInfo = Get-Item $OutputPath
            Write-Status "Built successfully: $([math]::Round($fileInfo.Length / 1MB, 1)) MB" "SUCCESS"
            
            # Quick test
            $version = & $OutputPath -version 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Status "Binary test passed: $version" "SUCCESS"
                return $true
            } else {
                Write-Status "Binary test failed" "ERROR"
                return $false
            }
        } else {
            Write-Status "Build failed: $result" "ERROR"
            return $false
        }
    }
    catch {
        Write-Status "Build error: $($_.Exception.Message)" "ERROR"
        return $false
    }
}

function Test-Templates {
    param([string]$BinaryPath, [string]$ProvidersDir, [string]$TargetProvider)
    
    Write-Status "Testing templates..." "TEST"
    
    if (-not (Test-Path $BinaryPath)) {
        Write-Status "Binary not found: $BinaryPath" "ERROR"
        return $false
    }
    
    if (-not (Test-Path $ProvidersDir)) {
        Write-Status "Providers directory not found: $ProvidersDir" "ERROR"
        return $false
    }
    
    $providers = Get-ChildItem -Path $ProvidersDir -Directory | Where-Object { Test-Path "$($_.FullName)\templates" }
    
    if ($TargetProvider) {
        $providers = $providers | Where-Object { $_.Name -eq $TargetProvider }
        if (-not $providers) {
            Write-Status "Provider not found: $TargetProvider" "ERROR"
            return $false
        }
    }
    
    $totalTests = 0
    $passedTests = 0
    
    foreach ($provider in $providers) {
        $templates = Get-ChildItem -Path "$($provider.FullName)\templates" -Filter "*.yaml" -Recurse
        
        if ($templates.Count -eq 0) {
            continue
        }
        
        Write-Status "Testing $($provider.Name): $($templates.Count) templates"
        
        foreach ($template in $templates) {
            $templateName = [System.IO.Path]::GetFileNameWithoutExtension($template.Name)
            $configPath = "$($provider.FullName)\configs\$($template.Name)"
            $resultPath = "$($provider.FullName)\results\$templateName.drawio"
            
            # Create config if needed
            if (-not (Test-Path $configPath)) {
                $configDir = Split-Path $configPath -Parent
                if (-not (Test-Path $configDir)) {
                    New-Item -ItemType Directory -Path $configDir -Force | Out-Null
                }
                
                # Get template name from file
                $actualTemplateName = $templateName
                try {
                    $templateContent = Get-Content -Path $template.FullName -Raw
                    if ($templateContent -match '(?m)^name:\s*["`'']?([^"`'']+)["`'']?') {
                        $actualTemplateName = $matches[1]
                    }
                } catch { }
                
                $testConfig = @"
version: "1.0"
diagram:
  pages:
    - id: "page1"
      name: "Test - $templateName"
      elements:
        - id: "test-id"
          name: "test-element"
          template: "$actualTemplateName"
          position:
            x: 100
            y: 100
          parameters:
"@
                
                # Add basic parameters
                try {
                    if ($templateContent -match "(?ms)^parameters:\s*\n(.*?)(?=^\S|\Z)") {
                        $parametersSection = $matches[1]
                        $paramLines = $parametersSection -split "`n" | Where-Object { $_ -match "^\s*-\s*name:" }
                        
                        foreach ($paramLine in $paramLines) {
                            if ($paramLine -match '^\s*-\s*name:\s*["`'']?([^"`'']+)["`'']?') {
                                $paramName = $matches[1]
                                $testConfig += "`n            $paramName`: `"TestValue`""
                            }
                        }
                    }
                } catch { }
                
                Set-Content -Path $configPath -Value $testConfig -Encoding UTF8
            }
            
            # Run test
            $resultDir = Split-Path $resultPath -Parent
            if (-not (Test-Path $resultDir)) {
                New-Item -ItemType Directory -Path $resultDir -Force | Out-Null
            }
            
            $templatesBasePath = "$($provider.FullName)\templates"
            $result = & $BinaryPath -input $configPath -output $resultPath -templates $templatesBasePath 2>&1
            $totalTests++
            
            if ($LASTEXITCODE -eq 0 -and (Test-Path $resultPath)) {
                try {
                    [xml]$xmlContent = Get-Content $resultPath
                    $passedTests++
                    if ($Verbose) { Write-Status "  ✓ $templateName" "SUCCESS" }
                }
                catch {
                    if ($Verbose) { Write-Status "  ✗ $templateName (invalid XML)" "ERROR" }
                }
            } else {
                if ($Verbose) { Write-Status "  ✗ $templateName (generation failed)" "ERROR" }
            }
        }
    }
    
    $failedTests = $totalTests - $passedTests
    Write-Status "Tests: $passedTests/$totalTests passed" $(if ($failedTests -eq 0) { "SUCCESS" } else { "WARNING" })
    
    return $failedTests -eq 0
}

# Main execution
if ($SkipTest) {
    Write-Host "Hippo - Build"
} elseif ($Build -and $Test) {
    Write-Host "Hippo - Build and Test"
} elseif ($Build -and -not $Test) {
    Write-Host "Hippo - Build"
} elseif (-not $Build -and $Test) {
    Write-Host "Hippo - Test"
} else {
    Write-Host "Hippo - Test"
}
Write-Host "======================"

$success = $true

if ($Build) {
    $buildSuccess = Build-Hippo -OutputPath $Binary -Clean $Clean -Verbose $Verbose
    if (-not $buildSuccess) {
        $success = $false
    }
}

if ($Test -and $success) {
    $testSuccess = Test-Templates -BinaryPath $Binary -ProvidersDir $ProvidersDir -TargetProvider $Provider
    if (-not $testSuccess) {
        $success = $false
    }
}

if ($success) {
    Write-Status "All operations completed successfully" "SUCCESS"
    exit 0
} else {
    Write-Status "Some operations failed" "ERROR"
    exit 1
}
