# KUBERNETES Provider Templates

This directory contains all templates, examples, and configurations for KUBERNETES cloud services.

## Directory Structure

- **templates/**: Template definitions for KUBERNETES services
- **examples/**: Example configurations and use cases
- **configs/**: Test configurations (auto-generated)
- **results/**: Generated DrawIO diagrams
- **docs/**: Provider-specific documentation

## Templates

### Workload Resources
- **deployment**: Kubernetes Deployments
- **pod**: Kubernetes Pods
- **service**: Kubernetes Services
- **ingress**: Ingress Controllers

### Configuration Resources
- **configmap**: Configuration Maps
- **secret**: Kubernetes Secrets
- **persistent-volume**: Persistent Volumes

### Cluster Resources
- **namespace**: Kubernetes Namespaces
- **cluster**: Kubernetes Clusters
- **node**: Kubernetes Nodes

## Usage

To use these templates in your diagrams:

```yaml
version: "1.0"
diagram:
  pages:
    - id: "page1"
      name: "KUBERNETES Architecture"
      elements:
        - name: "my-service"
          template: "kubernetes/[service-type]/[template-name].yaml"
          parameters:
            # Template-specific parameters
```

## Testing

Run tests for this provider:

```powershell
.\scripts\hippodamus-converter.ps1 -TestOnly -ProvidersDir "providers\kubernetes"
```

## Contributing

When adding new templates:

1. Place template files in the appropriate service subdirectory under 	emplates/
2. Add example configurations in xamples/
3. Update this README with template documentation
4. Ensure templates follow the unified group-based format
