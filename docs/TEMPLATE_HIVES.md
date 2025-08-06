# Template Hive System Implementation

## 🎯 Overview

The Template Hive System organizes templates into technology-specific namespaces, making it easier to manage, discover, and use templates for different platforms and domains.

## 🏗️ Architecture

### Hive Structure
```
templates/
├── azuredevops/          # Azure DevOps Platform Hive
│   ├── azuredevops-organization.yaml
│   ├── azuredevops-project.yaml
│   ├── azuredevops-repository.yaml
│   └── azuredevops-pipeline.yaml
├── kubernetes/           # Kubernetes Platform Hive  
│   ├── kubernetes-cluster.yaml
│   ├── kubernetes-namespace.yaml
│   └── kubernetes-pod.yaml
├── github/              # GitHub Platform Hive
│   ├── github-organization.yaml
│   └── github-repository.yaml
└── generic/             # Generic Infrastructure Hive
    ├── container.yaml
    ├── server.yaml
    ├── database.yaml
    ├── microservice.yaml
    └── loadbalancer.yaml
```

### Example Structure
```
examples/
├── azuredevops/         # Azure DevOps Examples
│   ├── azuredevops-simpler.yaml
│   ├── azuredevops-complete.yaml
│   └── test-*-deps.yaml
├── kubernetes/          # Kubernetes Examples
│   └── k8s-simple.yaml
└── generic/            # Generic Examples
    ├── infrastructure.yaml
    ├── microservices.yaml
    └── simple.yaml
```

## 🚀 Features

### 1. **Hive-Aware Template Loading**
- Automatically detects hive structure during template loading
- Registers templates with hive-qualified keys (`hive/template`)
- Maintains hive registry for discovery and validation

### 2. **Intelligent Template Resolution**
- **Within-hive resolution**: Templates automatically resolve to same hive
- **Cross-hive references**: Explicit `hive/template.yaml` syntax
- **Fallback resolution**: Falls back to root level templates

### 3. **Dependency Validation with Hives**
- Dependencies work across hive boundaries
- Validates parent relationships regardless of hive
- Supports mixed-hive configurations

### 4. **Hive Discovery API**
- `ListHives()` - Get all available hives
- `ListTemplatesInHive(hive)` - Get templates in specific hive
- `GetTemplateKey(name, hive)` - Generate hive-qualified keys

## 📝 Usage Examples

### Basic Hive Usage
```yaml
templates:
  - name: "k8s-cluster"
    template: "kubernetes/kubernetes-cluster.yaml"
  - name: "k8s-namespace"  
    template: "kubernetes/kubernetes-namespace.yaml"

diagram:
  pages:
    - elements:
        - template: "k8s-cluster"
          children:
            - template: "k8s-namespace"  # Auto-resolves to kubernetes/kubernetes-namespace
```

### Cross-Hive Configuration
```yaml
templates:
  - name: "azure-org"
    template: "azuredevops/azuredevops-organization.yaml"
  - name: "github-repo"
    template: "github/github-repository.yaml"
  - name: "k8s-cluster"
    template: "kubernetes/kubernetes-cluster.yaml"

# Mix templates from different platforms in same diagram
diagram:
  pages:
    - elements:
        - template: "azure-org"       # Azure DevOps hive
        - template: "github-repo"    # GitHub hive  
        - template: "k8s-cluster"    # Kubernetes hive
```

### Dependency Validation Across Hives
```yaml
# This works - namespace depends on cluster (both kubernetes hive)
cluster:
  template: "k8s-cluster"
  children:
    - template: "k8s-namespace"  ✅

# This fails - namespace without cluster parent
standalone:
  template: "k8s-namespace"  ❌ Error: requires kubernetes-cluster parent
```

## 🧪 Test Cases

### 1. **Azure DevOps Hive** ✅
```bash
go run ./cmd/hippodamus -i examples/azuredevops/azuredevops-simpler.yaml -t templates
# Result: Success - Organization → Project hierarchy validated
```

### 2. **Kubernetes Hive** ✅  
```bash
go run ./cmd/hippodamus -i examples/kubernetes/k8s-simple.yaml -t templates
# Result: Success - Cluster → Namespace → Pod hierarchy validated
```

### 3. **Cross-Hive Dependencies** ✅
```yaml
# Azure DevOps project with Kubernetes deployment
azuredevops-project:
  children:
    - template: "kubernetes/kubernetes-cluster"  # Cross-hive reference works
```

## 🎯 Benefits

### 1. **Organization & Discoverability**
- **Technology-specific grouping** makes templates easier to find
- **Clear namespace separation** prevents naming conflicts
- **Logical organization** by platform/domain

### 2. **Scalability**  
- **Easy to add new hives** for additional platforms (AWS, GCP, etc.)
- **Isolated template evolution** within each hive
- **Platform-specific best practices** can be encoded in each hive

### 3. **User Experience**
- **Intuitive structure** maps to real-world technology stacks
- **Auto-resolution** reduces configuration verbosity
- **Clear error messages** with hive context

### 4. **Maintainability**
- **Modular template organization** 
- **Technology-specific versioning** possible per hive
- **Easy template discovery** and maintenance

## 🔮 Future Hive Possibilities

### Cloud Provider Hives
```
templates/
├── aws/                 # Amazon Web Services
│   ├── aws-vpc.yaml
│   ├── aws-ec2.yaml
│   └── aws-rds.yaml
├── azure/              # Microsoft Azure
│   ├── azure-resource-group.yaml
│   ├── azure-vm.yaml
│   └── azure-sql.yaml
└── gcp/                # Google Cloud Platform
    ├── gcp-project.yaml
    ├── gcp-compute.yaml
    └── gcp-storage.yaml
```

### Technology Stack Hives
```
templates/
├── docker/             # Docker & Containers
├── terraform/          # Infrastructure as Code
├── jenkins/            # CI/CD Pipelines
├── monitoring/         # Observability Stack
└── networking/         # Network Infrastructure
```

## 📊 Implementation Summary

| Feature | Status | Description |
|---------|--------|-------------|
| **Hive Detection** | ✅ | Automatic hive detection from folder structure |
| **Template Loading** | ✅ | Hive-aware template loading and registration |
| **Template Resolution** | ✅ | Intelligent within-hive and cross-hive resolution |
| **Dependency Validation** | ✅ | Dependencies work across hive boundaries |
| **Discovery API** | ✅ | List hives and templates within hives |
| **Error Handling** | ✅ | Clear error messages with hive context |
| **Documentation** | ✅ | Comprehensive examples and usage patterns |

The Template Hive System provides a robust, scalable foundation for organizing diagram templates by technology platform, making Hippodamus suitable for complex, multi-technology architectures!
