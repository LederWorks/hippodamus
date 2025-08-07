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

### Changed
- Simplified resource syntax from verbose provider configuration to clean `resource: 'template-name'` format
- Provider declarations moved to top-level configuration
- Template workflow redesigned for better clarity and usability

### Technical Details
- Go modules support
- Cross-platform compatibility
- Automated testing and linting
- Semantic versioning with GitVersion
- GitHub Actions CI/CD pipeline
- Comprehensive repository automation

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
