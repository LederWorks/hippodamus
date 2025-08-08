#!/usr/bin/env pwsh
# Build script for Hippodamus with branch-based versioning
# Usage: .\build.ps1

param(
    [string]$OutputDir = "build",
    [string]$Version = "",
    [switch]$Clean = $false,
    [switch]$Verbose = $false
)

# Set error handling
$ErrorActionPreference = "Stop"

Write-Host "🔨 Hippodamus Build Script" -ForegroundColor Cyan
Write-Host "=========================" -ForegroundColor Cyan

# Get current branch name
try {
    $branch = git branch --show-current
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to get current branch"
    }
    Write-Host "📍 Current branch: $branch" -ForegroundColor Green
} catch {
    Write-Error "❌ Failed to get git branch information. Are you in a git repository?"
    exit 1
}

# Clean branch name for filename (replace special characters)
$cleanBranch = $branch -replace '[/\\:*?"<>|]', '-'

# Determine version
if ($Version -eq "") {
    # Try to get version from GitVersion if available
    try {
        $gitVersionOutput = gitversion 2>$null
        if ($LASTEXITCODE -eq 0) {
            $gitVersionJson = $gitVersionOutput | ConvertFrom-Json
            $Version = $gitVersionJson.SemVer
            Write-Host "📦 GitVersion detected: v$Version" -ForegroundColor Green
        } else {
            throw "GitVersion not available"
        }
    } catch {
        # Fallback to git describe or manual version
        try {
            $Version = git describe --tags --always --dirty 2>$null
            if ($LASTEXITCODE -ne 0) {
                $Version = "dev-$cleanBranch"
            }
            Write-Host "📦 Using git describe: $Version" -ForegroundColor Yellow
        } catch {
            $Version = "dev-$cleanBranch"
            Write-Host "📦 Using fallback version: $Version" -ForegroundColor Yellow
        }
    }
} else {
    Write-Host "📦 Using provided version: $Version" -ForegroundColor Green
}

# Get build timestamp and commit hash
$buildDate = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
try {
    $commit = git rev-parse HEAD
    if ($LASTEXITCODE -ne 0) {
        $commit = "unknown"
    }
} catch {
    $commit = "unknown"
}

# Get short commit hash (first 8 characters)
if ($commit -ne "unknown") {
    $commitShort = $commit.Substring(0, 8)
} else {
    $commitShort = "unknown"
}

# Build filename
$outputName = "hippodamus-$cleanBranch-$commitShort"
$outputPath = Join-Path $OutputDir $outputName

# Clean build directory if requested
if ($Clean -and (Test-Path $OutputDir)) {
    Write-Host "🧹 Cleaning build directory..." -ForegroundColor Yellow
    Remove-Item -Recurse -Force $OutputDir
}

# Create build directory
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
    Write-Host "📁 Created build directory: $OutputDir" -ForegroundColor Green
}

# Build information
Write-Host "🔧 Build Configuration:" -ForegroundColor Cyan
Write-Host "   Branch: $branch" -ForegroundColor White
Write-Host "   Version: $Version" -ForegroundColor White
Write-Host "   Commit: $commitShort" -ForegroundColor White
Write-Host "   Output: $outputPath.exe" -ForegroundColor White

# Build ldflags
$ldflags = "-X main.version=$Version -X main.commit=$commit -X main.date=$buildDate"

# Verbose output
if ($Verbose) {
    Write-Host "🔍 Verbose build information:" -ForegroundColor Cyan
    Write-Host "   GOOS: $env:GOOS" -ForegroundColor White
    Write-Host "   GOARCH: $env:GOARCH" -ForegroundColor White
    Write-Host "   LDFlags: $ldflags" -ForegroundColor White
    Write-Host "   Command: go build -ldflags `"$ldflags`" -o `"$outputPath.exe`" ./cmd/hippodamus" -ForegroundColor White
}

# Build the application
Write-Host "⚡ Building..." -ForegroundColor Yellow
try {
    $env:CGO_ENABLED = "0"  # Static binary
    
    if ($Verbose) {
        go build -v -ldflags $ldflags -o "$outputPath.exe" ./cmd/hippodamus
    } else {
        go build -ldflags $ldflags -o "$outputPath.exe" ./cmd/hippodamus
    }
    
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed with exit code $LASTEXITCODE"
    }
} catch {
    Write-Error "❌ Build failed: $_"
    exit 1
}

# Verify the build
if (Test-Path "$outputPath.exe") {
    $fileSize = (Get-Item "$outputPath.exe").Length
    $fileSizeMB = [math]::Round($fileSize / 1MB, 2)
    
    Write-Host "✅ Build successful!" -ForegroundColor Green
    Write-Host "📄 Output file: $outputPath.exe ($fileSizeMB MB)" -ForegroundColor Green
    
    # Test the built binary
    Write-Host "🧪 Testing built binary..." -ForegroundColor Cyan
    try {
        $versionOutput = & "$outputPath.exe" --version
        Write-Host "   Version output: $versionOutput" -ForegroundColor White
        
        $providersOutput = & "$outputPath.exe" --list-providers
        Write-Host "   Providers available: $(($providersOutput | Measure-Object -Line).Lines) lines of output" -ForegroundColor White
    } catch {
        Write-Warning "⚠️ Binary test failed, but build completed: $_"
    }
} else {
    Write-Error "❌ Build appeared to succeed but output file not found: $outputPath.exe"
    exit 1
}

Write-Host "`n🎉 Build complete! Run with:" -ForegroundColor Green
Write-Host "   .\$outputPath.exe --help" -ForegroundColor White
