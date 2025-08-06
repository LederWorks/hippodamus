# PowerShell script to find and convert remaining templates with "elements:" sections
param(
    [string]$TemplatesDir = "templates",
    [switch]$DryRun = $false
)

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

# Find templates that still need conversion
Write-Host "🔍 Finding templates that still have 'elements:' sections..."
Write-Host "Template directory: $TemplatesDir"
Write-Host ""

$templatesToConvert = @()

Get-ChildItem -Path $TemplatesDir -Recurse -Filter "*.yaml" | ForEach-Object { 
    $content = Get-Content $_.FullName -Raw
    if ($content -match "(?m)^elements:") { 
        $templatesToConvert += $_.FullName
        Write-Host "📄 Still needs conversion: $($_.FullName)"
    }
}

if ($templatesToConvert.Count -eq 0) {
    Write-Host "✅ All templates have been converted to the new group format!"
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
Get-ChildItem -Path $TemplatesDir -Recurse -Filter "*.yaml" | ForEach-Object { 
    if ((Get-Content $_.FullName -Raw) -match "(?m)^elements:") { 
        $remainingTemplates += $_.FullName
    }
}

if ($remainingTemplates.Count -eq 0) {
    Write-Host "✅ All templates successfully converted!"
} else {
    Write-Host "⚠️  $($remainingTemplates.Count) templates still need manual conversion:"
    $remainingTemplates | ForEach-Object { Write-Host "    $_" }
}
