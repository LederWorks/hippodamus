# Contributing to Hippodamus

Thank you for your interest in contributing to Hippodamus! We welcome contributions from the community and appreciate your help in making this project better.

## Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to [support@lederworks.com](mailto:support@lederworks.com).

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to avoid duplicates
2. **Use the issue templates** provided for bugs and feature requests
3. **Provide detailed information** including:
   - Go version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Error messages or logs

### Suggesting Features

We welcome feature suggestions! Please:

1. **Check existing feature requests** to avoid duplicates
2. **Use the feature request template**
3. **Explain the use case** and why it would be valuable
4. **Consider backward compatibility**

### Contributing Code

#### Prerequisites

- Go 1.21 or later
- Git
- golangci-lint (for code quality checks)

#### Development Setup

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/hippodamus.git
   cd hippodamus
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/LederWorks/hippodamus.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

#### Making Changes

1. **Create a feature branch** from main:
   ```bash
   git checkout -b feature/amazing-feature
   ```

2. **Make your changes** following these guidelines:
   - Write clear, self-documenting code
   - Add tests for new functionality
   - Update documentation as needed
   - Follow Go conventions and idioms

3. **Test your changes**:
   ```bash
   # Run tests
   go test ./...
   
   # Run linting
   golangci-lint run
   
   # Build the application
   go build -o hippodamus cmd/hippodamus/main.go
   ```

4. **Commit your changes** with clear commit messages:
   ```bash
   git add .
   git commit -m "feat: add amazing feature
   
   - Implement feature X
   - Add tests for feature X
   - Update documentation"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/amazing-feature
   ```

6. **Create a Pull Request** on GitHub

#### Pull Request Guidelines

- **Use the PR template** provided
- **Write a clear title** that summarizes the change
- **Describe the changes** in detail
- **Reference related issues** using keywords (fixes #123)
- **Ensure CI passes** before requesting review
- **Keep PRs focused** - one feature/fix per PR
- **Update CHANGELOG.md** if the change affects users

#### Code Style

We follow standard Go conventions:

- Use `gofmt` for formatting
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Write meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

#### Testing

- Write tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions
- Include integration tests for major features

#### Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update examples if the API changes
- Consider adding to the project wiki for complex features

## Project Structure

```
hippodamus/
├── cmd/hippodamus/          # Main application entry point
├── pkg/                     # Reusable packages
│   ├── drawio/             # Draw.io XML generation
│   ├── providers/          # Provider system
│   ├── schema/             # Configuration schema
│   └── templates/          # Template processing
├── examples/               # Example configurations
├── docs/                   # Documentation
├── .github/                # GitHub workflows and templates
└── tests/                  # Integration tests
```

## Release Process

Releases are handled automatically:

1. **Version Calculation**: GitVersion calculates the next version based on commit messages
2. **Automated Builds**: GitHub Actions builds for all platforms
3. **Automated Releases**: Releases are created automatically for version tags
4. **Changelog**: Updated automatically based on conventional commits

### Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning:

- `feat:` - New features (minor version bump)
- `fix:` - Bug fixes (patch version bump)
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

For breaking changes, add `BREAKING CHANGE:` in the commit body or use `!` after the type:
```
feat!: redesign provider API

BREAKING CHANGE: Provider interface has changed
```

## Community

- **Discussions**: GitHub Discussions for questions and ideas
- **Issues**: GitHub Issues for bugs and feature requests
- **Support**: Email [support@lederworks.com](mailto:support@lederworks.com)

## Recognition

Contributors will be recognized in:
- Release notes
- CHANGELOG.md
- GitHub contributors list

Thank you for contributing to Hippodamus! 🎉
