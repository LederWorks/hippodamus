# GCP Provider Templates

This directory contains all templates, examples, and configurations for GCP cloud services.

## Directory Structure

- **templates/**: Template definitions for GCP services
- **examples/**: Example configurations and use cases
- **configs/**: Test configurations (auto-generated)
- **results/**: Generated DrawIO diagrams
- **docs/**: Provider-specific documentation

## Templates

### Compute Services
- **compute-instance**: Google Compute Engine instances
- **gke-cluster**: Google Kubernetes Engine clusters
- **cloud-run**: Cloud Run services
- **cloud-functions**: Cloud Functions

### Storage Services
- **cloud-storage**: Cloud Storage buckets
- **persistent-disk**: Persistent Disks
- **filestore**: Cloud Filestore

### Networking Services
- **vpc-network**: Virtual Private Cloud networks
- **subnet**: VPC Subnets
- **firewall-rule**: Firewall Rules
- **load-balancer**: Cloud Load Balancers

### Management Services
- **organization**: GCP Organizations
- **project**: GCP Projects
- **folder**: Resource Folders
- **region**: GCP Regions

## Usage

To use these templates in your diagrams:

```yaml
version: "1.0"
diagram:
  pages:
    - id: "page1"
      name: "GCP Architecture"
      elements:
        - name: "my-service"
          template: "gcp/[service-type]/[template-name].yaml"
          parameters:
            # Template-specific parameters
```

## Testing

Run tests for this provider:

```powershell
.\scripts\hippodamus-converter.ps1 -TestOnly -ProvidersDir "providers\gcp"
```

## Contributing

When adding new templates:

1. Place template files in the appropriate service subdirectory under 	emplates/
2. Add example configurations in xamples/
3. Update this README with template documentation
4. Ensure templates follow the unified group-based format
