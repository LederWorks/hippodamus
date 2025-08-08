package core

import (
	"fmt"

	"github.com/LederWorks/hippodamus/pkg/providers"
	"github.com/LederWorks/hippodamus/pkg/schema"
	"github.com/LederWorks/hippodamus/providers/core/resources"
	"github.com/LederWorks/hippodamus/providers/core/templates"
)

// CoreProvider implements the Provider interface for core diagram elements
type CoreProvider struct {
	version string
	// Resource instances
	shapeResource     *resources.ShapeResource
	connectorResource *resources.ConnectorResource
	textResource      *resources.TextResource
	groupResource     *resources.GroupResource
	swimlaneResource  *resources.SwimlaneResource
	// Template instances
	shapeTemplate     *templates.ShapeTemplate
	connectorTemplate *templates.ConnectorTemplate
	textTemplate      *templates.TextTemplate
	groupTemplate     *templates.GroupTemplate
	swimlaneTemplate  *templates.SwimlaneTemplate
}

// NewCoreProvider creates a new core provider instance
func NewCoreProvider() *CoreProvider {
	return &CoreProvider{
		version:           "1.0.0",
		shapeResource:     resources.NewShapeResource(),
		connectorResource: resources.NewConnectorResource(),
		textResource:      resources.NewTextResource(),
		groupResource:     resources.NewGroupResource(),
		swimlaneResource:  resources.NewSwimlaneResource(),
		shapeTemplate:     templates.NewShapeTemplate(),
		connectorTemplate: templates.NewConnectorTemplate(),
		textTemplate:      templates.NewTextTemplate(),
		groupTemplate:     templates.NewGroupTemplate(),
		swimlaneTemplate:  templates.NewSwimlaneTemplate(),
	}
}

// Name returns the provider name
func (p *CoreProvider) Name() string {
	return "core"
}

// Version returns the provider version
func (p *CoreProvider) Version() string {
	return p.version
}

// Resources returns the list of supported core elements
func (p *CoreProvider) Resources() []providers.ResourceDefinition {
	return []providers.ResourceDefinition{
		p.shapeResource.Definition(),
		p.connectorResource.Definition(),
		p.textResource.Definition(),
		p.groupResource.Definition(),
		p.swimlaneResource.Definition(),
	}
}

// Validate validates core element parameters
func (p *CoreProvider) Validate(resourceType string, params map[string]interface{}) error {
	switch resourceType {
	case "shape":
		return p.shapeResource.Validate(params)
	case "connector":
		return p.connectorResource.Validate(params)
	case "text":
		return p.textResource.Validate(params)
	case "group":
		return p.groupResource.Validate(params)
	case "swimlane":
		return p.swimlaneResource.Validate(params)
	default:
		return &providers.ProviderError{
			Provider: p.Name(),
			Resource: resourceType,
			Message:  fmt.Sprintf("unsupported resource type: %s", resourceType),
			Code:     "UNSUPPORTED_RESOURCE",
		}
	}
}

// GenerateTemplate generates core element templates
func (p *CoreProvider) GenerateTemplate(resourceType string, params map[string]interface{}) (*schema.Element, error) {
	if err := p.Validate(resourceType, params); err != nil {
		return nil, err
	}

	switch resourceType {
	case "shape":
		return p.shapeTemplate.Generate(params)
	case "connector":
		return p.connectorTemplate.Generate(params)
	case "text":
		return p.textTemplate.Generate(params)
	case "group":
		return p.groupTemplate.Generate(params)
	case "swimlane":
		return p.swimlaneTemplate.Generate(params)
	default:
		return nil, &providers.ProviderError{
			Provider: p.Name(),
			Resource: resourceType,
			Message:  fmt.Sprintf("unsupported resource type: %s", resourceType),
			Code:     "UNSUPPORTED_RESOURCE",
		}
	}
}

// GetSchema returns the JSON schema for a resource type
func (p *CoreProvider) GetSchema(resourceType string) (map[string]interface{}, error) {
	resources := p.Resources()
	for _, resource := range resources {
		if resource.Type == resourceType {
			return resource.Schema, nil
		}
	}

	return nil, &providers.ProviderError{
		Provider: p.Name(),
		Resource: resourceType,
		Message:  fmt.Sprintf("schema not found for resource type: %s", resourceType),
		Code:     "SCHEMA_NOT_FOUND",
	}
}
