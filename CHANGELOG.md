# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Hippodamus
- YAML to Draw.io XML converter functionality
- Clean provider/template syntax system
- Support for three provider types:
  - `builtin`: Local provider folders (exact name matching)
  - `registry`: LederWorks GitHub organization providers (default)
  - `custom`: Filesystem or HTTPS Git repository providers
- Template hive system for bulk template loading
- Template collections with include/exclude patterns
- Provider addressing with namespace resolution
- GitVersion integration for automated semantic versioning
- Comprehensive CI/CD pipeline with GitHub Actions
- Multi-platform builds (Linux, Windows, macOS)
- Automated dependency management with Dependabot
- Code quality enforcement with golangci-lint
- Automated releases with changelog generation
- Dynamic provider version injection from GitVersion/build system
- Core provider system with modular architecture
- Built-in provider registry with automatic initialization
- Cross-platform build scripts (`build.ps1` and `build.bat`)
- Branch-based build naming for development builds
- Build artifacts now include commit ID in filename format: `hippodamus-{branch}-{commit8}`
- Comprehensive branch protection documentation (`docs/BRANCH_PROTECTION.md`)
- Team permission matrices for repository governance
- Skip-release mechanism for development commits
- GitHub organization and repository configuration:
  - Created and configured teams at organization level
  - Implemented branch protection rules on repository
  - Set up team-based access controls and permissions
  - Configured repository settings for automated workflows

### Changed
- Simplified resource syntax from verbose provider configuration to clean `resource: 'template-name'` format
- Provider declarations moved to top-level configuration
- Template workflow redesigned for better clarity and usability
- Updated GitHub Actions to latest versions (v4/v5/v6) to resolve deprecation warnings
- Core provider version now dynamically set from application version instead of hardcoded value
- Provider initialization moved to main application startup for proper version injection
- CODEOWNERS simplified to use individual maintainer instead of non-existent teams
- Core provider test file renamed from `provider_modular_test.go` to `provider_test.go`
- Test function names cleaned up, removed "Modular" prefix for better readability
- Updated README.md documentation to reflect current implementation:
  - Provider Types section updated to show builtin providers vs planned registry/custom providers
  - Quick Start section updated with build script usage and realistic YAML examples
  - Template system examples updated to reflect current core provider capabilities
  - Added Contributing section with development workflow including build scripts
  - Added License section referencing MIT License file
  - Removed references to non-existent AWS templates and unrealistic configuration examples
  - Updated all configuration examples to use current builtin provider system
- Build scripts updated to include commit ID in output filename format: `hippodamus-{branch}-{commit8}`

### Fixed
- GitHub Actions deprecation warnings in CI/CD pipeline
- Provider version synchronization with GitVersion releases
- Core provider registration in global registry
- Branch protection configuration for team-based workflows
- Duplicate CI pipeline runs caused by overlapping push and pull_request triggers
- golangci-lint configuration compatibility with version 1.64.8 (removed deprecated 'version' field)
- Code formatting issues across all Go files to meet linting standards

### Security
- Updated all GitHub Actions to latest secure versions
- Implemented proper team permission structures for repository access
- Configured GitHub organization teams with appropriate access levels
- Applied branch protection rules and team-based review requirements

### Technical Details
- Go modules support
- Cross-platform compatibility
- Automated testing and linting
- Semantic versioning with GitVersion
- GitHub Actions CI/CD pipeline
- Comprehensive repository automation
- Provider system architecture with dynamic registration
- Build-time version injection using Go ldflags
- Multi-platform build automation with PowerShell and Batch scripts
- Repository governance with branch protection and team permissions
- GitHub organization and repository administration:
  - Organization-level team creation and management
  - Repository-level branch protection rule implementation
  - Team permission matrix configuration for secure workflows
  - Automated workflow integration with organizational policies

## [0.1.0] - TBD

### Added
- Initial public release
- Core YAML to Draw.io conversion functionality
- Provider system implementation
- Template hive support
- Documentation and examples

---

## Release Notes

### v0.1.0 - Initial Release

This is the first public release of Hippodamus, a powerful YAML to Draw.io XML converter that simplifies diagram creation through templates and providers.

#### Key Features:
- **Clean Syntax**: Simple `resource: 'template-name'` format
- **Provider System**: Support for builtin, registry, and custom providers
- **Template Hives**: Bulk template loading with pattern matching
- **Automation**: Full CI/CD pipeline with automated versioning and releases

#### Getting Started:
1. Download the latest release for your platform
2. Create a YAML configuration file
3. Run `hippodamus input.yaml output.xml`
4. Open the generated XML file in Draw.io

For detailed documentation and examples, see the [README](README.md).
