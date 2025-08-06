# Hippodamus - YAML to D# Use templates
./hippodamus -input diagram.yaml -templates templates/ -output diagram.drawio

# Auto-detect output format based on input filename
./hippodamus -input diagram.yaml
```

## Documentation

- **[Development Guide](docs/DEVELOPMENT.md)** - Development setup, building, and contributing
- **[Template Dependencies](docs/TEMPLATE_DEPENDENCIES.md)** - Understanding template dependency system
- **[Dependency Implementation](docs/DEPENDENCY_IMPLEMENTATION.md)** - Technical details of dependency validation
- **[Template Hives](docs/TEMPLATE_HIVES.md)** - Template organization by technology platform
- **[Configuration Comparison](docs/CONFIG_COMPARISON.md)** - Simple vs. explicit configuration approaches
- **[Reorganization Summary](docs/REORGANIZATION_SUMMARY.md)** - Recent project structure improvements XML Converter

Hippodamus is a Go application that transforms YAML configuration files into draw.io XML syntax, allowing you to programmatically create diagrams from structured configuration files.

## Features

- Transform YAML configurations to draw.io XML format
- Support for nested diagrams with parents, layers, and tags
- Template system for reusable diagram components
- Customizable YAML schema for different diagram types
- Support for complex diagram structures including shapes, connectors, and styling

## Usage

```bash
# Build the application
go build -o hippodamus ./cmd/hippodamus

# Convert YAML to draw.io format (default .drawio extension)
./hippodamus -input diagram.yaml -output diagram.drawio

# Convert to XML format
./hippodamus -input diagram.yaml -output diagram.xml

# Use templates
./hippodamus -input diagram.yaml -templates templates/ -output diagram.drawio

# Auto-detect output format based on input filename
./hippodamus -input diagram.yaml
```

## Supported Output Formats

- `.drawio` - Draw.io native format (default)
- `.xml` - Standard XML format

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
