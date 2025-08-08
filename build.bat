@echo off
REM Build script for Hippodamus with branch-based versioning
REM Usage: build.bat

setlocal enabledelayedexpansion

echo 🔨 Hippodamus Build Script
echo =========================

REM Get current branch name
for /f "tokens=*" %%i in ('git branch --show-current 2^>nul') do set BRANCH=%%i
if errorlevel 1 (
    echo ❌ Failed to get git branch information. Are you in a git repository?
    exit /b 1
)

echo 📍 Current branch: %BRANCH%

REM Clean branch name for filename (replace special characters)
set CLEAN_BRANCH=%BRANCH%
set CLEAN_BRANCH=%CLEAN_BRANCH:/=-%
set CLEAN_BRANCH=%CLEAN_BRANCH:\=-%

REM Try to get version from GitVersion
gitversion >nul 2>&1
if %errorlevel% equ 0 (
    for /f "tokens=2 delims=:" %%i in ('gitversion ^| findstr "SemVer"') do (
        set VERSION_RAW=%%i
        set VERSION=!VERSION_RAW: "=!
        set VERSION=!VERSION:",=!
        set VERSION=!VERSION:"=!
    )
    echo 📦 GitVersion detected: v!VERSION!
) else (
    REM Fallback to git describe
    for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
    if errorlevel 1 (
        set VERSION=dev-%CLEAN_BRANCH%
    )
    echo 📦 Using fallback version: !VERSION!
)

REM Get build timestamp and commit hash
for /f "tokens=*" %%i in ('powershell -command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'"') do set BUILD_DATE=%%i
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT_SHORT=%%i
if errorlevel 1 (
    set COMMIT_SHORT=unknown
) else (
    REM Also get full commit for ldflags
    for /f "tokens=*" %%i in ('git rev-parse HEAD 2^>nul') do set COMMIT=%%i
    if errorlevel 1 set COMMIT=unknown
)

if "%COMMIT_SHORT%"=="unknown" set COMMIT=unknown

REM Build information
set OUTPUT_NAME=hippodamus-%CLEAN_BRANCH%-%COMMIT_SHORT%
set OUTPUT_DIR=build
set OUTPUT_PATH=%OUTPUT_DIR%\%OUTPUT_NAME%.exe

echo 🔧 Build Configuration:
echo    Branch: %BRANCH%
echo    Version: !VERSION!
echo    Commit: %COMMIT_SHORT%
echo    Output: %OUTPUT_PATH%

REM Create build directory
if not exist "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%"
    echo 📁 Created build directory: %OUTPUT_DIR%
)

REM Build ldflags
set LDFLAGS=-X main.version=!VERSION! -X main.commit=%COMMIT% -X main.date=%BUILD_DATE%

REM Build the application
echo ⚡ Building...
set CGO_ENABLED=0
go build -ldflags "%LDFLAGS%" -o "%OUTPUT_PATH%" ./cmd/hippodamus

if errorlevel 1 (
    echo ❌ Build failed
    exit /b 1
)

REM Verify the build
if exist "%OUTPUT_PATH%" (
    echo ✅ Build successful!
    echo 📄 Output file: %OUTPUT_PATH%
    
    REM Test the built binary
    echo 🧪 Testing built binary...
    "%OUTPUT_PATH%" --version
    if errorlevel 1 (
        echo ⚠️ Binary test failed, but build completed
    ) else (
        echo    Binary test passed
    )
) else (
    echo ❌ Build appeared to succeed but output file not found: %OUTPUT_PATH%
    exit /b 1
)

echo.
echo 🎉 Build complete! Run with:
echo    %OUTPUT_PATH% --help

endlocal
