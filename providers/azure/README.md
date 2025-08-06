# AZURE Provider Templates

This directory contains all templates, examples, and configurations for AZURE cloud services.

## Directory Structure

- **templates/**: Template definitions for AZURE services
- **examples/**: Example configurations and use cases
- **configs/**: Test configurations (auto-generated)
- **results/**: Generated DrawIO diagrams
- **docs/**: Provider-specific documentation

## Templates

### Compute Services
- **virtual-machine**: Azure Virtual Machines
- **aks-cluster**: Azure Kubernetes Service clusters
- **container-instances**: Azure Container Instances
- **function-app**: Azure Functions

### Storage Services
- **storage-account**: Azure Storage Accounts
- **blob-storage**: Blob Storage containers
- **file-share**: Azure File Shares

### Networking Services
- **virtual-network**: Azure Virtual Networks
- **subnet**: VNet Subnets
- **network-security-group**: Network Security Groups
- **load-balancer**: Azure Load Balancers

### Management Services
- **management-group**: Azure Management Groups
- **subscription**: Azure Subscriptions
- **resource-group**: Resource Groups
- **tenant**: Azure AD Tenants

## Usage

To use these templates in your diagrams:

```yaml
version: "1.0"
diagram:
  pages:
    - id: "page1"
      name: "AZURE Architecture"
      elements:
        - name: "my-service"
          template: "azure/[service-type]/[template-name].yaml"
          parameters:
            # Template-specific parameters
```

## Testing

Run tests for this provider:

```powershell
.\scripts\hippodamus-converter.ps1 -TestOnly -ProvidersDir "providers\azure"
```

## Contributing

When adding new templates:

1. Place template files in the appropriate service subdirectory under 	emplates/
2. Add example configurations in xamples/
3. Update this README with template documentation
4. Ensure templates follow the unified group-based format
