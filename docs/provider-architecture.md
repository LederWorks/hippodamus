# Hippodamus Provider Architecture Design

## Overview
This document outlines the evolution from static YAML providers to a Terraform-inspired plugin architecture for Hippodamus providers.

## Current State vs. Target Architecture

### Phase 1: Current (YAML-based)
```
hippodamus/
├── providers/
│   ├── aws/
│   │   ├── configs/
│   │   └── templates/
│   ├── azure/
│   └── gcp/
```

### Phase 2: Hybrid (Go + YAML)
```
hippodamus/
├── pkg/
│   └── providers/
│       ├── interface.go      # Provider interface
│       ├── registry.go       # Provider registry
│       └── loader.go         # Dynamic loading
├── providers/
│   ├── aws/
│   │   ├── provider.go       # Go implementation
│   │   ├── resources.go      # Resource definitions
│   │   ├── templates.go      # Template generation
│   │   ├── provider_test.go  # Provider tests
│   │   └── configs/          # Example configs
│   ├── azure/
│   └── gcp/
```

### Phase 3: Plugin Architecture (External)
```
hippodamus-provider-aws/     # Separate repository
├── main.go                  # Plugin entrypoint
├── provider/
│   ├── resources.go
│   ├── templates.go
│   └── validation.go
├── go.mod                   # Independent versioning
└── examples/

hippodamus/
├── pkg/
│   └── providers/
│       ├── plugin/          # Plugin system
│       ├── grpc/           # gRPC communication
│       └── discovery/      # Provider discovery
```

## Implementation Phases

### Phase 1: Provider Interface Design
1. **Define Provider Interface**
   ```go
   type Provider interface {
       Name() string
       Version() string
       Resources() []ResourceDefinition
       GenerateTemplate(resource string, params map[string]interface{}) (*schema.Element, error)
       Validate(config map[string]interface{}) error
   }
   ```

2. **Create Provider Registry**
   - Central registration system
   - Version management
   - Dependency resolution

3. **Implement Built-in Providers**
   - Convert existing AWS/Azure/GCP YAML to Go
   - Maintain backward compatibility
   - Add comprehensive testing

### Phase 2: Plugin System
1. **gRPC-based Communication**
   - Define provider protocol
   - Implement plugin loading
   - Handle provider lifecycle

2. **Provider SDK**
   - Helper libraries for provider authors
   - Testing frameworks
   - Documentation generation

3. **Discovery Mechanism**
   - Local provider discovery
   - Remote provider registries
   - Version resolution

### Phase 3: Ecosystem Development
1. **External Provider Repositories**
   - Separate GitHub repositories
   - Independent CI/CD
   - Community contributions

2. **Provider Marketplace**
   - Provider registry
   - Documentation hub
   - Version management

## Benefits

### For Core Development
- **Focused scope**: Core team focuses on diagram engine
- **Faster releases**: Provider updates don't require core releases
- **Better testing**: Isolated provider testing
- **Reduced complexity**: Smaller, more maintainable codebase

### For Provider Authors
- **Independent development**: Own repositories and release cycles
- **Rich tooling**: Go's ecosystem for validation, testing, documentation
- **Type safety**: Compile-time validation of provider logic
- **Flexible distribution**: Multiple ways to distribute providers

### For Users
- **Selective installation**: Install only needed providers
- **Version flexibility**: Use different provider versions per project
- **Better documentation**: Provider-specific docs and examples
- **Community ecosystem**: Access to community-contributed providers

## Technical Considerations

### Backward Compatibility
- Maintain YAML template support during transition
- Automatic migration tools for existing configurations
- Clear deprecation timeline

### Performance
- Plugin caching to avoid startup overhead
- Efficient gRPC communication
- Lazy loading of providers

### Security
- Provider signature verification
- Sandboxed execution environment
- Permission-based resource access

## Migration Strategy

### Week 1-2: Foundation
- Implement provider interface
- Create provider registry
- Convert one provider (AWS) to Go

### Week 3-4: Validation
- Add comprehensive testing
- Performance benchmarking
- Community feedback

### Week 5-6: Plugin System
- Implement gRPC protocol
- Create provider SDK
- External provider proof-of-concept

### Week 7-8: Ecosystem
- Documentation and examples
- Community onboarding
- Release preparation

## Success Metrics

- **Provider development velocity**: Time to create new providers
- **Community adoption**: Number of external providers
- **Performance**: Plugin loading and execution times
- **User experience**: Migration ease and functionality
