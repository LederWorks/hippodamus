# GENERIC Provider Templates

This directory contains all templates, examples, and configurations for GENERIC cloud services.

## Directory Structure

- **templates/**: Template definitions for GENERIC services
- **examples/**: Example configurations and use cases
- **configs/**: Test configurations (auto-generated)
- **results/**: Generated DrawIO diagrams
- **docs/**: Provider-specific documentation

## Templates

### Infrastructure Components
- **container**: Generic container/pod representation
- **database**: Database services
- **load-balancer**: Load balancing components
- **microservice**: Microservice applications
- **server**: Generic server instances

## Usage

To use these templates in your diagrams:

```yaml
version: "1.0"
diagram:
  pages:
    - id: "page1"
      name: "GENERIC Architecture"
      elements:
        - name: "my-service"
          template: "generic/[service-type]/[template-name].yaml"
          parameters:
            # Template-specific parameters
```

## Testing

Run tests for this provider:

```powershell
.\scripts\hippodamus-converter.ps1 -TestOnly -ProvidersDir "providers\generic"
```

## Contributing

When adding new templates:

1. Place template files in the appropriate service subdirectory under 	emplates/
2. Add example configurations in xamples/
3. Update this README with template documentation
4. Ensure templates follow the unified group-based format
