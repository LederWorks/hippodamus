# AWS Provider Templates

This directory contains all templates, examples, and configurations for AWS cloud services.

## Directory Structure

- **templates/**: Template definitions for AWS services
- **examples/**: Example configurations and use cases
- **configs/**: Test configurations (auto-generated)
- **results/**: Generated DrawIO diagrams
- **docs/**: Provider-specific documentation

## Templates

### Compute Services
- **ec2-instance**: EC2 virtual machines
- **ecs-cluster**: Elastic Container Service clusters
- **eks-cluster**: Elastic Kubernetes Service clusters
- **lambda-function**: Serverless Lambda functions

### Storage Services
- **s3-bucket**: Simple Storage Service buckets
- **ebs-volume**: Elastic Block Store volumes
- **efs-filesystem**: Elastic File System

### Networking Services
- **vpc**: Virtual Private Cloud
- **subnet**: VPC Subnets
- **security-group**: Security Groups
- **load-balancer**: Application/Network Load Balancers

### Management Services
- **organization**: AWS Organizations
- **account**: AWS Accounts
- **organizational-unit**: Organizational Units
- **region**: AWS Regions

## Usage

To use these templates in your diagrams:

```yaml
version: "1.0"
diagram:
  pages:
    - id: "page1"
      name: "AWS Architecture"
      elements:
        - name: "my-service"
          template: "aws/[service-type]/[template-name].yaml"
          parameters:
            # Template-specific parameters
```

## Testing

Run tests for this provider:

```powershell
.\scripts\hippodamus-converter.ps1 -TestOnly -ProvidersDir "providers\aws"
```

## Contributing

When adding new templates:

1. Place template files in the appropriate service subdirectory under 	emplates/
2. Add example configurations in xamples/
3. Update this README with template documentation
4. Ensure templates follow the unified group-based format
