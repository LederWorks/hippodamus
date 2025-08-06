# PowerShell script to convert templates from Elements format to Group format
param(
    [string]$TemplatesDir = "templates"
)

function Convert-Template {
    param([string]$FilePath)
    
    Write-Host "Converting $FilePath..."
    
    $content = Get-Content -Path $FilePath -Raw
    
    # Skip if already has 'group:' field
    if ($content -match "(?m)^group:") {
        Write-Host "  Already converted, skipping."
        return
    }
    
    # Skip if no 'elements:' field
    if (-not ($content -match "(?m)^elements:")) {
        Write-Host "  No elements field found, skipping."
        return
    }
    
    try {
        # Parse YAML-like structure using regex
        # Extract the main container element and its children
        
        # Find elements section
        if ($content -match "(?ms)^elements:\s*\n(.*?)(?=^\S|\Z)") {
            $elementsSection = $matches[1]
            
            # Extract first element (main container)
            if ($elementsSection -match "(?ms)^\s*-\s*type:\s*[`"']?container[`"']?\s*\n(.*?)(?=^\s*-\s*type:|\Z)") {
                $mainElement = $matches[1]
                
                # Extract properties, style, nesting, and children from main element
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
                    $newContent = $newContent -replace "(?m)^dependencies:", "dependencies:`n$customDep"
                }
                
                # Write the converted content
                Set-Content -Path $FilePath -Value $newContent -NoNewline
                Write-Host "  Converted successfully."
            }
        }
    }
    catch {
        Write-Error "Error converting $FilePath`: $($_.Exception.Message)"
    }
}

# Find all YAML templates
$templates = Get-ChildItem -Path $TemplatesDir -Recurse -Filter "*.yaml"

Write-Host "Found $($templates.Count) template files to process..."

foreach ($template in $templates) {
    Convert-Template -FilePath $template.FullName
}

Write-Host "Template conversion completed!"
