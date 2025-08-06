# PowerShell script to convert templates and validate with functional testing
param(
    [string]$ProvidersDir = "providers",
    [string]$TemplatesDir = "",  # Legacy support - will be ignored if ProvidersDir exists
    [string]$ConfigDir = "",     # Legacy support - will be ignored if ProvidersDir exists
    [string]$ResultsDir = "",    # Legacy support - will be ignored if ProvidersDir exists
    [string]$HippodamusExecutable = ".\hippo.exe",
    [string]$Provider = "",      # Specific provider to test (optional)
    [switch]$DryRun = $false,
    [switch]$SkipTesting = $false,
    [switch]$TestOnly = $false,
    [switch]$LegacyMode = $false  # Force use of old structure
)

function Get-StructureInfo {
    param()
    
    # Determine if we're using the new provider structure or legacy structure
    if ((Test-Path $ProvidersDir) -and -not $LegacyMode) {
        return @{
            Mode = "Provider"
            ProvidersDir = $ProvidersDir
            Providers = Get-ChildItem -Path $ProvidersDir -Directory | Where-Object { Test-Path "$($_.FullName)\templates" }
        }
    } else {
        # Legacy mode - use old structure
        $legacyTemplatesDir = if ($TemplatesDir) { $TemplatesDir } else { "templates" }
        $legacyConfigDir = if ($ConfigDir) { $ConfigDir } else { "config" }
        $legacyResultsDir = if ($ResultsDir) { $ResultsDir } else { "results" }
        
        return @{
            Mode = "Legacy"
            TemplatesDir = $legacyTemplatesDir
            ConfigDir = $legacyConfigDir
            ResultsDir = $legacyResultsDir
        }
    }
}

function Get-ProviderTemplates {
    param(
        [string]$ProviderPath
    )
    
    $templatesPath = Join-Path $ProviderPath "templates"
    if (Test-Path $templatesPath) {
        return Get-ChildItem -Path $templatesPath -Filter "*.yaml" -Recurse
    }
    return @()
}

function Test-TemplateWithHippodamus {
    param(
        [string]$TemplatePath,
        [string]$ConfigPath,
        [string]$OutputPath,
        [string]$HippodamusExecutable,
        [string]$TemplatesBasePath = "templates"
    )
    
    Write-Host "🧪 Testing template: $(Split-Path $TemplatePath -Leaf)..."
    
    # Ensure output directory exists
    $outputDir = Split-Path $OutputPath -Parent
    if (-not (Test-Path $outputDir)) {
        New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
    }
    
    try {
        # Run hippodamus with the config file
        # Use the provider-specific templates path for the -templates parameter
        $templatesParam = if ($TemplatesBasePath -and (Test-Path $TemplatesBasePath)) { $TemplatesBasePath } else { "templates" }
        $result = & $HippodamusExecutable -input $ConfigPath -output $OutputPath -templates $templatesParam 2>&1
        $exitCode = $LASTEXITCODE
        
        if ($exitCode -eq 0 -and (Test-Path $OutputPath)) {
            Write-Host "  ✅ Test passed - DrawIO file generated successfully"
            
            # Verify the output file is valid XML
            try {
                [xml]$xmlContent = Get-Content $OutputPath
                Write-Host "  ✅ Generated DrawIO file is valid XML"
                return $true
            }
            catch {
                Write-Host "  ❌ Generated file is not valid XML: $($_.Exception.Message)"
                return $false
            }
        }
        else {
            Write-Host "  ❌ Test failed with exit code: $exitCode"
            Write-Host "  Output: $result"
            return $false
        }
    }
    catch {
        Write-Host "  ❌ Error running hippodamus: $($_.Exception.Message)"
        return $false
    }
}

function Create-TestConfigForTemplate {
    param(
        [string]$TemplatePath,
        [string]$ConfigOutputPath,
        [string]$TemplateBasePath = ""
    )
    
    # Ensure config directory exists
    $configDir = Split-Path $ConfigOutputPath -Parent
    if (-not (Test-Path $configDir)) {
        New-Item -ItemType Directory -Path $configDir -Force | Out-Null
    }
    
    # Extract template name and relative path
    $templateName = [System.IO.Path]::GetFileNameWithoutExtension($TemplatePath)
    
    # Calculate relative path for template reference
    # Read the actual template name from the template file
    $actualTemplateName = $relativePath
    try {
        $templateContent = Get-Content -Path $TemplatePath -Raw
        if ($templateContent -match '(?m)^name:\s*["`'']?([^"`'']+)["`'']?') {
            $actualTemplateName = $matches[1]
        }
    } catch {
        # If we can't read the template, fall back to the calculated name
    }
    
    $templateReference = if ($TemplateBasePath) {
        # For provider-first structure, use the actual template name from the file
        $actualTemplateName
    } else {
        # For legacy mode
        [System.IO.Path]::GetRelativePath("templates", $TemplatePath) -replace "\\", "/"
    }
    
    # Create a test configuration with proper structure
    $testConfig = @"
version: "1.0"
diagram:
  pages:
    - id: "page1"
      name: "Test Page - $templateName"
      elements:
        - id: "$templateName-test-id"
          name: "$templateName-test"
          template: "$templateReference"
          position:
            x: 100
            y: 100
          parameters:
"@
    
    # Try to read the template to extract parameter information
    try {
        $templateContent = Get-Content -Path $TemplatePath -Raw
        
        # Extract parameters from template
        if ($templateContent -match "(?ms)^parameters:\s*\n(.*?)(?=^\S|\Z)") {
            $parametersSection = $matches[1]
            
            # Parse each parameter and generate test values
            $paramLines = $parametersSection -split "`n" | Where-Object { $_ -match "^\s*-\s*name:" }
            
            foreach ($paramLine in $paramLines) {
                if ($paramLine -match '^\s*-\s*name:\s*["`'']?([^"`'']+)["`'']?') {
                    $paramName = $matches[1]
                    
                    # Look for default value or type to generate appropriate test value
                    $paramBlock = ""
                    if ($parametersSection -match "(?ms)name:\s*[`"']?$paramName[`"']?\s*\n(.*?)(?=^\s*-\s*name:|\Z)") {
                        $paramBlock = $matches[1]
                    }
                    
                    $testValue = ""
                    if ($paramBlock -match 'default:\s*["`'']?([^"`'']+)["`'']?') {
                        $testValue = $matches[1]
                    }
                    elseif ($paramBlock -match 'type:\s*["`'']?string["`'']?') {
                        $testValue = "Test$paramName"
                    }
                    elseif ($paramBlock -match 'type:\s*["`'']?color["`'']?') {
                        $testValue = "#FF0000"
                    }
                    elseif ($paramBlock -match 'type:\s*["`'']?number["`'']?') {
                        $testValue = "100"
                    }
                    else {
                        $testValue = "TestValue"
                    }
                    
                    $testConfig += "`n            $paramName`: `"$testValue`""
                }
            }
        }
    }
    catch {
        Write-Host "  ⚠️  Could not parse template parameters, using basic config"
    }
    
    # Write the test configuration
    Set-Content -Path $ConfigOutputPath -Value $testConfig -Encoding UTF8
}

function Test-AllTemplates {
    param(
        [object]$StructureInfo,
        [string]$HippodamusExecutable,
        [string]$SpecificProvider = ""
    )
    
    Write-Host ""
    Write-Host "🧪 Starting comprehensive template testing..."
    Write-Host "Structure Mode: $($StructureInfo.Mode)"
    Write-Host "Executable: $HippodamusExecutable"
    Write-Host ""
    
    # Check if hippodamus executable exists
    if (-not (Test-Path $HippodamusExecutable)) {
        Write-Error "❌ Hippodamus executable not found at: $HippodamusExecutable"
        Write-Host "Please build the project first: go build -o hippo.exe ./cmd/hippodamus"
        return $false
    }
    
    $allTests = @()
    $passedTests = 0
    $failedTests = 0
    
    if ($StructureInfo.Mode -eq "Provider") {
        # New provider-based structure
        $providersToTest = if ($SpecificProvider) {
            $StructureInfo.Providers | Where-Object { $_.Name -eq $SpecificProvider }
        } else {
            $StructureInfo.Providers
        }
        
        foreach ($provider in $providersToTest) {
            Write-Host "🔧 Testing provider: $($provider.Name)"
            
            $templates = Get-ProviderTemplates -ProviderPath $provider.FullName
            $templatesBasePath = Join-Path $provider.FullName "templates"
            
            foreach ($template in $templates) {
                $relativePath = [System.IO.Path]::GetRelativePath($templatesBasePath, $template.FullName)
                $configPath = Join-Path $provider.FullName "configs" ($relativePath -replace "\.yaml$", ".yaml")
                $resultPath = Join-Path $provider.FullName "results" ($relativePath -replace "\.yaml$", ".drawio")
                
                # Create test config if it doesn't exist
                if (-not (Test-Path $configPath)) {
                    Write-Host "📝 Creating test config: $configPath"
                    Create-TestConfigForTemplate -TemplatePath $template.FullName -ConfigOutputPath $configPath -TemplateBasePath $templatesBasePath
                }
                
                # Run the test with provider-specific template path
                $testResult = Test-TemplateWithHippodamus -TemplatePath $template.FullName -ConfigPath $configPath -OutputPath $resultPath -HippodamusExecutable $HippodamusExecutable -TemplatesBasePath $templatesBasePath
                
                $testInfo = @{
                    Provider = $provider.Name
                    Template = $relativePath
                    ConfigPath = $configPath
                    ResultPath = $resultPath
                    Passed = $testResult
                }
                $allTests += $testInfo
                
                if ($testResult) {
                    $passedTests++
                } else {
                    $failedTests++
                }
            }
        }
    } else {
        # Legacy structure
        Write-Host "Templates: $($StructureInfo.TemplatesDir)"
        Write-Host "Configs: $($StructureInfo.ConfigDir)"
        Write-Host "Results: $($StructureInfo.ResultsDir)"
        
        $templates = Get-ChildItem -Path $StructureInfo.TemplatesDir -Recurse -Filter "*.yaml"
        
        foreach ($template in $templates) {
            $relativePath = [System.IO.Path]::GetRelativePath($StructureInfo.TemplatesDir, $template.FullName)
            $configPath = Join-Path $StructureInfo.ConfigDir ($relativePath -replace "\.yaml$", ".yaml")
            $resultPath = Join-Path $StructureInfo.ResultsDir ($relativePath -replace "\.yaml$", ".drawio")
            
            # Create test config if it doesn't exist
            if (-not (Test-Path $configPath)) {
                Write-Host "📝 Creating test config: $configPath"
                Create-TestConfigForTemplate -TemplatePath $template.FullName -ConfigOutputPath $configPath
            }
            
            # Run the test
            $testResult = Test-TemplateWithHippodamus -TemplatePath $template.FullName -ConfigPath $configPath -OutputPath $resultPath -HippodamusExecutable $HippodamusExecutable
            
            $testInfo = @{
                Template = $relativePath
                ConfigPath = $configPath
                ResultPath = $resultPath
                Passed = $testResult
            }
            $allTests += $testInfo
            
            if ($testResult) {
                $passedTests++
            } else {
                $failedTests++
            }
        }
    }
    
    # Print summary
    Write-Host ""
    Write-Host "🎯 Testing Summary:"
    Write-Host "  Total templates tested: $($allTests.Count)"
    Write-Host "  ✅ Passed: $passedTests"
    Write-Host "  ❌ Failed: $failedTests"
    Write-Host ""
    
    if ($failedTests -gt 0) {
        Write-Host "Failed tests:"
        $allTests | Where-Object { -not $_.Passed } | ForEach-Object {
            $displayName = if ($_.Provider) { "$($_.Provider)/$($_.Template)" } else { $_.Template }
            Write-Host "  ❌ $displayName"
        }
    }
    
    return $failedTests -eq 0
}

function Convert-ShapeTemplate {
    param(
        [string]$FilePath,
        [bool]$DryRun = $false
    )
    
    Write-Host "Converting shape template: $FilePath..."
    
    if ($DryRun) {
        Write-Host "  [DRY RUN] Would convert this template"
        return
    }
    
    $content = Get-Content -Path $FilePath -Raw
    
    try {
        # Extract elements section
        if ($content -match "(?ms)^elements:\s*\n(.*?)(?=^\S|\Z)") {
            $elementsSection = $matches[1]
            
            # Extract the first (and likely only) element
            if ($elementsSection -match "(?ms)^\s*-\s*type:\s*[`"']?shape[`"']?\s*\n(.*?)(?=^\s*-\s*type:|\Z)") {
                $shapeElement = $matches[1]
                
                # Extract properties and style
                $properties = ""
                $style = ""
                
                if ($shapeElement -match "(?ms)^\s*properties:\s*\n(.*?)(?=^\s*\w+:|\Z)") {
                    $properties = $matches[1]
                }
                
                if ($shapeElement -match "(?ms)^\s*style:\s*\n(.*?)(?=^\s*\w+:|\Z)") {
                    $style = $matches[1]
                }
                
                # Build new group section
                $groupSection = "group:`n"
                if ($properties) {
                    $groupSection += "  properties:`n" + ($properties -replace "(?m)^    ", "    ")
                }
                if ($style) {
                    $groupSection += "  style:`n" + ($style -replace "(?m)^    ", "    ")
                }
                
                # Add default container behavior
                $groupSection += "  autoResize: false`n"
                $groupSection += "  arrangement: `"free`"`n"
                
                # Replace elements section with group section
                $newContent = $content -replace "(?ms)^elements:.*?(?=^\S|\Z)", $groupSection
                
                # Add custom parent dependency if not present
                if (-not ($newContent -match "(?m)^\s*-\s*name:\s*[`"']?custom[`"']?")) {
                    $customDep = @"
  - name: "custom"
    type: "container" #container or group, defaults to container
    required: false
    relationship: "parent"
    description: "Custom parent object"
"@
                    
                    if ($newContent -match "(?m)^dependencies:") {
                        $newContent = $newContent -replace "(?m)^dependencies:", "dependencies:`n$customDep"
                    } else {
                        # Add dependencies section after version
                        $newContent = $newContent -replace "(?m)^version:.*?`n", "$&`ndependencies:`n$customDep`n"
                    }
                }
                
                # Write the converted content
                Set-Content -Path $FilePath -Value $newContent -NoNewline
                Write-Host "  ✅ Converted successfully."
            } else {
                Write-Host "  ⚠️  No shape element found in elements section"
            }
        } else {
            Write-Host "  ⚠️  No elements section found"
        }
    }
    catch {
        Write-Error "  ❌ Error converting $FilePath`: $($_.Exception.Message)"
    }
}

function Convert-ContainerTemplate {
    param(
        [string]$FilePath,
        [bool]$DryRun = $false
    )
    
    Write-Host "Converting container template: $FilePath..."
    
    if ($DryRun) {
        Write-Host "  [DRY RUN] Would convert this template"
        return
    }
    
    $content = Get-Content -Path $FilePath -Raw
    
    try {
        # This handles templates with container elements (more complex)
        # Extract elements section
        if ($content -match "(?ms)^elements:\s*\n(.*?)(?=^\S|\Z)") {
            $elementsSection = $matches[1]
            
            # Extract first container element
            if ($elementsSection -match "(?ms)^\s*-\s*type:\s*[`"']?container[`"']?\s*\n(.*?)(?=^\s*-\s*type:|\Z)") {
                $mainElement = $matches[1]
                
                # Extract sections
                $properties = ""
                $style = ""
                $nesting = ""
                $children = ""
                
                if ($mainElement -match "(?ms)^\s*properties:\s*\n(.*?)(?=^\s*\w+:|\Z)") {
                    $properties = $matches[1]
                }
                
                if ($mainElement -match "(?ms)^\s*style:\s*\n(.*?)(?=^\s*\w+:|\Z)") {
                    $style = $matches[1]
                }
                
                if ($mainElement -match "(?ms)^\s*nesting:\s*\n(.*?)(?=^\s*\w+:|\Z)") {
                    $nesting = $matches[1]
                }
                
                if ($mainElement -match "(?ms)^\s*children:\s*\n(.*?)(?=^\s*\w+:|\Z)") {
                    $children = $matches[1]
                }
                
                # Build new group section
                $groupSection = "group:`n"
                if ($properties) {
                    $groupSection += "  properties:`n" + ($properties -replace "(?m)^    ", "    ")
                }
                if ($style) {
                    $groupSection += "  style:`n" + ($style -replace "(?m)^    ", "    ")
                }
                
                # Extract nesting properties
                if ($nesting) {
                    if ($nesting -match "(?m)^\s*autoResize:\s*(.*)") {
                        $groupSection += "  autoResize: $($matches[1])`n"
                    }
                    if ($nesting -match "(?ms)^\s*padding:\s*\n(.*?)(?=^\s*\w+:|\Z)") {
                        $groupSection += "  padding:`n" + ($matches[1] -replace "(?m)^    ", "    ")
                    }
                    if ($nesting -match "(?m)^\s*spacing:\s*(.*)") {
                        $groupSection += "  spacing: $($matches[1])`n"
                    }
                    if ($nesting -match "(?m)^\s*arrangement:\s*(.*)") {
                        $groupSection += "  arrangement: $($matches[1])`n"
                    }
                } else {
                    # Add defaults
                    $groupSection += "  autoResize: true`n"
                    $groupSection += "  arrangement: `"free`"`n"
                }
                
                if ($children) {
                    $groupSection += "  children:`n" + ($children -replace "(?m)^    ", "    ")
                }
                
                # Replace elements section with group section
                $newContent = $content -replace "(?ms)^elements:.*?(?=^\S|\Z)", $groupSection
                
                # Add custom parent dependency if not present
                if (-not ($newContent -match "(?m)^\s*-\s*name:\s*[`"']?custom[`"']?")) {
                    $customDep = @"
  - name: "custom"
    type: "container" #container or group, defaults to container
    required: false
    relationship: "parent"
    description: "Custom parent object"
"@
                    
                    if ($newContent -match "(?m)^dependencies:") {
                        $newContent = $newContent -replace "(?m)^dependencies:", "dependencies:`n$customDep"
                    } else {
                        $newContent = $newContent -replace "(?m)^version:.*?`n", "$&`ndependencies:`n$customDep`n"
                    }
                }
                
                # Write the converted content
                Set-Content -Path $FilePath -Value $newContent -NoNewline
                Write-Host "  ✅ Converted successfully."
            } else {
                Write-Host "  ⚠️  No container element found, trying shape conversion..."
                Convert-ShapeTemplate -FilePath $FilePath -DryRun $DryRun
            }
        }
    }
    catch {
        Write-Error "  ❌ Error converting $FilePath`: $($_.Exception.Message)"
    }
}

# Main execution logic
Write-Host "🚀 Hippodamus Template Converter & Tester"
Write-Host "==========================================="
Write-Host ""

# Determine structure mode
$structureInfo = Get-StructureInfo

Write-Host "📁 Structure Mode: $($structureInfo.Mode)"
if ($structureInfo.Mode -eq "Provider") {
    Write-Host "📁 Providers Directory: $($structureInfo.ProvidersDir)"
    Write-Host "📁 Available Providers: $($structureInfo.Providers.Name -join ', ')"
} else {
    Write-Host "📁 Templates Directory: $($structureInfo.TemplatesDir)"
    Write-Host "📁 Config Directory: $($structureInfo.ConfigDir)"
    Write-Host "📁 Results Directory: $($structureInfo.ResultsDir)"
}

if ($Provider) {
    Write-Host "🎯 Target Provider: $Provider"
}
Write-Host ""

if ($TestOnly) {
    Write-Host "🧪 Running in TEST-ONLY mode - no template conversion"
    $testSuccess = Test-AllTemplates -StructureInfo $structureInfo -HippodamusExecutable $HippodamusExecutable -SpecificProvider $Provider
    if ($testSuccess) {
        Write-Host "🎉 All tests passed!"
        exit 0
    } else {
        Write-Host "❌ Some tests failed!"
        exit 1
    }
}

# Find templates that still need conversion
Write-Host "🔍 Finding templates that still have 'elements:' sections..."
Write-Host ""

$templatesToConvert = @()

if ($structureInfo.Mode -eq "Provider") {
    # Search in provider structure
    $providersToCheck = if ($Provider) {
        $structureInfo.Providers | Where-Object { $_.Name -eq $Provider }
    } else {
        $structureInfo.Providers
    }
    
    foreach ($providerDir in $providersToCheck) {
        $templates = Get-ProviderTemplates -ProviderPath $providerDir.FullName
        foreach ($template in $templates) {
            $content = Get-Content $template.FullName -Raw
            if ($content -match "(?m)^elements:") { 
                $templatesToConvert += $template.FullName
                $relPath = [System.IO.Path]::GetRelativePath($structureInfo.ProvidersDir, $template.FullName)
                Write-Host "📄 Still needs conversion: $relPath"
            }
        }
    }
} else {
    # Legacy structure
    Get-ChildItem -Path $structureInfo.TemplatesDir -Recurse -Filter "*.yaml" | ForEach-Object { 
        $content = Get-Content $_.FullName -Raw
        if ($content -match "(?m)^elements:") { 
            $templatesToConvert += $_.FullName
            Write-Host "📄 Still needs conversion: $($_.FullName)"
        }
    }
}

if ($templatesToConvert.Count -eq 0) {
    Write-Host "✅ All templates have been converted to the new group format!"
    
    if (-not $SkipTesting) {
        Write-Host ""
        Write-Host "🧪 Starting functional testing since all templates are converted..."
        $testSuccess = Test-AllTemplates -StructureInfo $structureInfo -HippodamusExecutable $HippodamusExecutable -SpecificProvider $Provider
        if ($testSuccess) {
            Write-Host "🎉 All conversions validated successfully!"
            exit 0
        } else {
            Write-Host "❌ Some templates failed validation!"
            exit 1
        }
    }
    exit 0
}

Write-Host ""
Write-Host "Found $($templatesToConvert.Count) templates that need conversion."

if ($DryRun) {
    Write-Host ""
    Write-Host "🔍 DRY RUN MODE - No files will be modified"
    Write-Host ""
}

foreach ($templatePath in $templatesToConvert) {
    Write-Host ""
    
    $content = Get-Content -Path $templatePath -Raw
    
    # Determine conversion strategy based on content
    if ($content -match "(?ms)^\s*-\s*type:\s*[`"']?container[`"']?") {
        Convert-ContainerTemplate -FilePath $templatePath -DryRun $DryRun
    } elseif ($content -match "(?ms)^\s*-\s*type:\s*[`"']?shape[`"']?") {
        Convert-ShapeTemplate -FilePath $templatePath -DryRun $DryRun
    } else {
        Write-Host "⚠️  Unknown template structure in: $templatePath"
        Write-Host "    Please convert manually."
    }
}

Write-Host ""
Write-Host "🎉 Template conversion completed!"

# Show final status
$remainingTemplates = @()

if ($structureInfo.Mode -eq "Provider") {
    # Search in provider structure
    $providersToCheck = if ($Provider) {
        $structureInfo.Providers | Where-Object { $_.Name -eq $Provider }
    } else {
        $structureInfo.Providers
    }
    
    foreach ($providerDir in $providersToCheck) {
        $templates = Get-ProviderTemplates -ProviderPath $providerDir.FullName
        foreach ($template in $templates) {
            if ((Get-Content $template.FullName -Raw) -match "(?m)^elements:") { 
                $remainingTemplates += $template.FullName
            }
        }
    }
} else {
    # Legacy structure
    Get-ChildItem -Path $structureInfo.TemplatesDir -Recurse -Filter "*.yaml" | ForEach-Object { 
        if ((Get-Content $_.FullName -Raw) -match "(?m)^elements:") { 
            $remainingTemplates += $_.FullName
        }
    }
}

if ($remainingTemplates.Count -eq 0) {
    Write-Host "✅ All templates successfully converted!"
    
    if (-not $SkipTesting -and -not $DryRun) {
        Write-Host ""
        Write-Host "🧪 Starting functional testing since all templates are now converted..."
        $testSuccess = Test-AllTemplates -StructureInfo $structureInfo -HippodamusExecutable $HippodamusExecutable -SpecificProvider $Provider
        if ($testSuccess) {
            Write-Host "🎉 All conversions validated successfully!"
            exit 0
        } else {
            Write-Host "❌ Some templates failed validation!"
            exit 1
        }
    }
} else {
    Write-Host "⚠️  $($remainingTemplates.Count) templates still need manual conversion:"
    $remainingTemplates | ForEach-Object { 
        if ($structureInfo.Mode -eq "Provider") {
            $relPath = [System.IO.Path]::GetRelativePath($structureInfo.ProvidersDir, $_)
            Write-Host "    $relPath"
        } else {
            Write-Host "    $_"
        }
    }
}
