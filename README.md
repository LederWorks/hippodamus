# Hippodamus

[![CI](https://github.com/LederWorks/hippodamus/actions/workflows/ci.yml/badge.svg)](https://github.com/LederWorks/hippodamus/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/LederWorks/hippodamus)](https://github.com/LederWorks/hippodamus/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/LederWorks/hippodamus)](https://goreportcard.com/report/github.com/LederWorks/hippodamus)
[![License](https://img.shields.io/github/license/LederWorks/hippodamus)](LICENSE)

> 🎯 **Hippodamus** - A powerful YAML to Draw.io XML converter that simplifies diagram creation through templates and providers.

Hippodamus transforms your YAML configuration files into Draw.io compatible XML diagrams using a flexible template and provider system. Perfect for automating diagram generation in CI/CD pipelines, documentation workflows, and infrastructure-as-code projects.

## ✨ Features

- **🎨 Clean Syntax**: Simple `resource: 'template-name'` format
- **🔌 Provider System**: Support for builtin, registry, and custom providers
- **📦 Template Hives**: Bulk template loading with pattern matching
- **🌍 Cross-Platform**: Works on Linux, Windows, and macOS
- **🚀 CI/CD Ready**: Perfect for automated documentation workflows
- **📊 Flexible Templates**: Extensible template system for any diagram type

## 🚀 Quick Start

### Installation

Download the latest release for your platform from the [releases page](https://github.com/LederWorks/hippodamus/releases).

### Basic Usage

1. Create a YAML configuration file:

```yaml
version: "1.0"
metadata:
  title: "My Infrastructure Diagram"
  description: "AWS infrastructure overview"

providers:
  aws:
    type: registry
    source: LederWorks/aws-templates

resources:
  - resource: 'vpc'
    name: 'main-vpc'
    cidr: '10.0.0.0/16'
  
  - resource: 'ec2-instance'
    name: 'web-server'
    instance_type: 't3.micro'
```

2. Convert to Draw.io XML:

```bash
hippodamus config.yaml diagram.xml
```

3. Open `diagram.xml` in Draw.io

## 📖 Documentation

### Provider Types

Hippodamus supports three types of providers:

#### 1. Builtin Providers
Local provider folders with exact name matching:
```yaml
providers:
  core:
    type: builtin
```

#### 2. Registry Providers (Default)
Templates from the LederWorks GitHub organization:
```yaml
providers:
  aws:
    type: registry
    source: LederWorks/aws-templates
```

#### 3. Custom Providers
Filesystem or HTTPS Git repository providers:
```yaml
providers:
  custom:
    type: custom
    source: "https://github.com/company/templates.git"
```

### Template Hives

Load multiple templates at once with pattern matching:

```yaml
template_hives:
  - name: "aws-core"
    source: "templates/aws"
    include:
      - "*.yaml"
      - "networking/**"
    exclude:
      - "*test*"
      - "deprecated/**"
```

### Configuration Reference

#### Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | string | ✅ | Configuration version |
| `metadata` | object | ✅ | Diagram metadata |
| `providers` | object | ❌ | Provider definitions |
| `template_hives` | array | ❌ | Template hive definitions |
| `resources` | array | ✅ | Resource definitions |

#### Provider Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✅ | Provider type: `builtin`, `registry`, `custom` |
| `source` | string | ❌ | Source location (for registry/custom types) |

#### Resource Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `resource` | string | ✅ | Template name to use |
| `name` | string | ❌ | Resource instance name |
| Additional fields depend on the specific template |

## 🏗️ Examples

### AWS Infrastructure

```yaml
version: "1.0"
metadata:
  title: "AWS Three-Tier Architecture"
  description: "Web application infrastructure"

providers:
  aws:
    type: registry
    source: LederWorks/aws-templates

resources:
  - resource: 'vpc'
    name: 'production-vpc'
    cidr: '10.0.0.0/16'
    
  - resource: 'subnet'
    name: 'public-subnet-1'
    vpc: 'production-vpc'
    cidr: '10.0.1.0/24'
    availability_zone: 'us-east-1a'
    
  - resource: 'application-load-balancer'
    name: 'web-alb'
    scheme: 'internet-facing'
    
  - resource: 'ec2-instance'
    name: 'web-server-1'
    instance_type: 't3.medium'
    subnet: 'public-subnet-1'
```

### Kubernetes Deployment

```yaml
version: "1.0"
metadata:
  title: "Kubernetes Application"
  description: "Microservices deployment"

providers:
  k8s:
    type: registry
    source: LederWorks/kubernetes-templates

resources:
  - resource: 'namespace'
    name: 'production'
    
  - resource: 'deployment'
    name: 'api-server'
    namespace: 'production'
    replicas: 3
    image: 'myapp/api:v1.2.0'
    
  - resource: 'service'
    name: 'api-service'
    namespace: 'production'
    type: 'ClusterIP'
    selector:
      app: 'api-server'
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or later
- Git

### Building from Source

```bash
git clone https://github.com/LederWorks/hippodamus.git
cd hippodamus
go build -o hippodamus cmd/hippodamus/main.go
```

### Running Tests

```bash
go test ./...
```

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Run linting: `golangci-lint run`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎯 Roadmap

- [ ] GUI application for visual diagram editing
- [ ] More built-in template providers
- [ ] Plugin system for custom renderers
- [ ] Integration with popular diagramming tools
- [ ] Template marketplace
- [ ] Real-time collaboration features

## 💬 Support

- 📖 [Documentation](https://github.com/LederWorks/hippodamus/wiki)
- 🐛 [Issue Tracker](https://github.com/LederWorks/hippodamus/issues)
- 💡 [Feature Requests](https://github.com/LederWorks/hippodamus/issues/new?template=feature_request.md)
- 📧 [Contact](mailto:support@lederworks.com)

## 🙏 Acknowledgments

- [Draw.io](https://www.draw.io/) for the excellent diagramming platform
- [GitVersion](https://gitversion.net/) for semantic versioning
- The Go community for amazing tools and libraries

---

**Made with ❤️ by [LederWorks](https://github.com/LederWorks)**

Both formats contain the same XML structure and are fully compatible with draw.io/diagrams.net.

## YAML Schema

The application supports a comprehensive YAML schema for defining draw.io diagrams:

```yaml
version: "1.0"
metadata:
  title: "Sample Diagram"
  description: "A sample diagram configuration"
  
templates:
  - name: "server"
    template: "templates/server.yaml"
    
diagram:
  pages:
    - id: "page1"
      name: "Main Page"
      layers:
        - id: "layer1"
          name: "Infrastructure"
          elements:
            - type: "shape"
              id: "server1"
              template: "server"
              properties:
                x: 100
                y: 100
                width: 120
                height: 80
                label: "Web Server"
```

## Project Structure

- `cmd/hippodamus/` - Main application entry point
- `pkg/schema/` - YAML schema definitions and Go structs
- `pkg/drawio/` - Draw.io XML generation logic
- `pkg/templates/` - Template processing system
- `templates/` - Reusable diagram templates
  - `azuredevops/` - Azure DevOps specific templates
- `examples/` - Example YAML configurations
- `docs/` - Documentation files
- `schemas/` - JSON schemas for YAML validation

## Templates

The template system supports reusable diagram components with dependency validation organized into **template hives** for different technology platforms:

### Template Hives

Template hives are organized folders that group related templates by technology platform or domain:

#### **Azure DevOps Hive** (`templates/azuredevops/`)
Azure DevOps resources with proper hierarchical dependencies:
- **Organization** (`azuredevops-organization.yaml`) - Root container
- **Project** (`azuredevops-project.yaml`) - Requires organization parent
- **Repository** (`azuredevops-repository.yaml`) - Requires project parent
- **Pipeline** (`azuredevops-pipeline.yaml`) - Requires repository parent
- **Environment, Library, Service Connection** - Various Azure DevOps resources

#### **Kubernetes Hive** (`templates/kubernetes/`)
Kubernetes resources with cluster hierarchy:
- **Cluster** (`kubernetes-cluster.yaml`) - Root Kubernetes cluster
- **Namespace** (`kubernetes-namespace.yaml`) - Requires cluster parent
- **Pod** (`kubernetes-pod.yaml`) - Requires namespace parent

#### **GitHub Hive** (`templates/github/`)
GitHub platform resources:
- **Organization** (`github-organization.yaml`) - GitHub organization
- **Repository** (`github-repository.yaml`) - Can be standalone or under organization

#### **Generic Hive** (`templates/generic/`)
General-purpose infrastructure components:
- `container.yaml` - Generic container shape
- `server.yaml` - Server infrastructure component
- `database.yaml` - Database component
- `microservice.yaml` - Microservice component
- `loadbalancer.yaml` - Load balancer component

### Template Reference Syntax

Templates can be referenced using hive notation:

```yaml
templates:
  - name: "k8s-cluster"
    template: "kubernetes/kubernetes-cluster.yaml"  # Hive-qualified reference
  - name: "github-repo"
    template: "github/github-repository.yaml"       # Cross-hive reference
```

### Template Resolution

The system supports intelligent template resolution:
- **Hive-aware resolution**: Templates automatically resolve within their hive context
- **Cross-hive references**: Use `hive/template.yaml` syntax for explicit references
- **Fallback resolution**: Falls back to root level if template not found in current hive

## Example Configurations

Examples are organized by technology hive in the `examples/` directory:

### Azure DevOps Examples (`examples/azuredevops/`)
- `azuredevops-simpler.yaml` - Minimal Azure DevOps configuration using template defaults
- `azuredevops-complete.yaml` - Comprehensive Azure DevOps setup with explicit configuration
- `test-*-deps.yaml` - Dependency validation test cases

### Kubernetes Examples (`examples/kubernetes/`)
- `k8s-simple.yaml` - Simple Kubernetes cluster with namespaces and pods

### Generic Examples (`examples/generic/`)
- `infrastructure.yaml` - Infrastructure diagram example  
- `microservices.yaml` - Microservices architecture example
- `simple.yaml` - Basic diagram example

### Usage Examples

#### Azure DevOps Organization
```yaml
templates:
  - name: "azuredevops-organization"
    template: "azuredevops/azuredevops-organization.yaml"

diagram:
  pages:
    - elements:
        - template: "azuredevops-organization"
          children:
            - template: "azuredevops-project"  # Auto-resolves within hive
```

#### Multi-Hive Configuration
```yaml
templates:
  - name: "k8s-cluster"
    template: "kubernetes/kubernetes-cluster.yaml"
  - name: "github-repo"
    template: "github/github-repository.yaml"

# Mix templates from different hives in the same diagram
```
