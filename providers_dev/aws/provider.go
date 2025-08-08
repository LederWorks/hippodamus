package aws

import (
	"fmt"

	"github.com/LederWorks/hippodamus/pkg/providers"
	"github.com/LederWorks/hippodamus/pkg/schema"
)

// AWSProvider implements the Provider interface for AWS resources
type AWSProvider struct {
	version string
}

// NewAWSProvider creates a new AWS provider instance
func NewAWSProvider() *AWSProvider {
	return &AWSProvider{
		version: "1.0.0",
	}
}

// Name returns the provider name
func (p *AWSProvider) Name() string {
	return "aws"
}

// Version returns the provider version
func (p *AWSProvider) Version() string {
	return p.version
}

// Resources returns the list of supported AWS resources
func (p *AWSProvider) Resources() []providers.ResourceDefinition {
	return []providers.ResourceDefinition{
		{
			Type:        "organization",
			Name:        "AWS Organization",
			Description: "AWS Organizations management structure",
			Category:    "management",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"orgName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the AWS Organization",
					},
					"managementAccountId": map[string]interface{}{
						"type":        "string",
						"description": "Management account ID",
						"pattern":     "^[0-9]{12}$",
					},
					"fillColor": map[string]interface{}{
						"type":        "string",
						"description": "Background color",
						"default":     "#FFF8E1",
					},
					"strokeColor": map[string]interface{}{
						"type":        "string",
						"description": "Border color",
						"default":     "#FF9900",
					},
				},
				"required": []string{"orgName", "managementAccountId"},
			},
			Examples: []providers.ResourceExample{
				{
					Name:        "Basic Organization",
					Description: "Simple AWS Organization setup",
					Config: map[string]interface{}{
						"orgName":             "My AWS Organization",
						"managementAccountId": "123456789012",
						"fillColor":           "#FFF8E1",
						"strokeColor":         "#FF9900",
					},
				},
			},
		},
		{
			Type:        "vpc",
			Name:        "Virtual Private Cloud",
			Description: "AWS VPC network infrastructure",
			Category:    "network",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"vpcName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the VPC",
					},
					"cidrBlock": map[string]interface{}{
						"type":        "string",
						"description": "CIDR block for the VPC",
						"pattern":     "^([0-9]{1,3}\\.){3}[0-9]{1,3}/[0-9]{1,2}$",
					},
					"region": map[string]interface{}{
						"type":        "string",
						"description": "AWS region",
						"default":     "us-east-1",
					},
				},
				"required": []string{"vpcName", "cidrBlock"},
			},
			Examples: []providers.ResourceExample{
				{
					Name:        "Basic VPC",
					Description: "Simple VPC configuration",
					Config: map[string]interface{}{
						"vpcName":   "production-vpc",
						"cidrBlock": "10.0.0.0/16",
						"region":    "us-east-1",
					},
				},
			},
		},
	}
}

// Validate validates AWS resource parameters
func (p *AWSProvider) Validate(resourceType string, params map[string]interface{}) error {
	switch resourceType {
	case "organization":
		return p.validateOrganization(params)
	case "vpc":
		return p.validateVPC(params)
	default:
		return &providers.ProviderError{
			Provider: p.Name(),
			Resource: resourceType,
			Message:  fmt.Sprintf("unsupported resource type: %s", resourceType),
			Code:     "UNSUPPORTED_RESOURCE",
		}
	}
}

// GenerateTemplate generates AWS resource templates
func (p *AWSProvider) GenerateTemplate(resourceType string, params map[string]interface{}) (*schema.Element, error) {
	if err := p.Validate(resourceType, params); err != nil {
		return nil, err
	}

	switch resourceType {
	case "organization":
		return p.generateOrganizationTemplate(params)
	case "vpc":
		return p.generateVPCTemplate(params)
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
func (p *AWSProvider) GetSchema(resourceType string) (map[string]interface{}, error) {
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

// validateOrganization validates AWS Organization parameters
func (p *AWSProvider) validateOrganization(params map[string]interface{}) error {
	if orgName, ok := params["orgName"]; !ok || orgName == "" {
		return &providers.ValidationError{
			Field:   "orgName",
			Message: "orgName is required",
			Code:    "REQUIRED_FIELD",
		}
	}

	if accountId, ok := params["managementAccountId"]; !ok || accountId == "" {
		return &providers.ValidationError{
			Field:   "managementAccountId",
			Message: "managementAccountId is required",
			Code:    "REQUIRED_FIELD",
		}
	}

	return nil
}

// validateVPC validates VPC parameters
func (p *AWSProvider) validateVPC(params map[string]interface{}) error {
	if vpcName, ok := params["vpcName"]; !ok || vpcName == "" {
		return &providers.ValidationError{
			Field:   "vpcName",
			Message: "vpcName is required",
			Code:    "REQUIRED_FIELD",
		}
	}

	if cidrBlock, ok := params["cidrBlock"]; !ok || cidrBlock == "" {
		return &providers.ValidationError{
			Field:   "cidrBlock",
			Message: "cidrBlock is required",
			Code:    "REQUIRED_FIELD",
		}
	}

	return nil
}

// generateOrganizationTemplate generates an AWS Organization template
func (p *AWSProvider) generateOrganizationTemplate(params map[string]interface{}) (*schema.Element, error) {
	orgName := params["orgName"].(string)
	fillColor := p.getStringParam(params, "fillColor", "#FFF8E1")
	strokeColor := p.getStringParam(params, "strokeColor", "#FF9900")

	return &schema.Element{
		Type: schema.ElementTypeShape,
		ID:   "aws-organization",
		Name: "AWS Organization",
		Properties: schema.ElementProperties{
			X:      100,
			Y:      100,
			Width:  300,
			Height: 200,
			Label:  orgName,
			Shape:  "rectangle",
		},
		Style: schema.Style{
			FillColor:     fillColor,
			StrokeColor:   strokeColor,
			StrokeWidth:   2,
			Rounded:       true,
			FontSize:      14,
			FontStyle:     "bold",
			TextAlign:     "center",
			VerticalAlign: "top",
		},
		Nesting: schema.NestingConfig{
			Mode:        schema.NestingModeChild,
			AutoResize:  true,
			Spacing:     15,
			Arrangement: schema.ArrangementVertical,
			Padding: schema.Padding{
				Top:    25,
				Right:  15,
				Bottom: 15,
				Left:   15,
			},
		},
	}, nil
}

// generateVPCTemplate generates a VPC template
func (p *AWSProvider) generateVPCTemplate(params map[string]interface{}) (*schema.Element, error) {
	vpcName := params["vpcName"].(string)
	cidrBlock := params["cidrBlock"].(string)
	region := p.getStringParam(params, "region", "us-east-1")

	label := fmt.Sprintf("%s\\n%s\\n%s", vpcName, cidrBlock, region)

	return &schema.Element{
		Type: schema.ElementTypeShape,
		ID:   "aws-vpc",
		Name: "AWS VPC",
		Properties: schema.ElementProperties{
			X:      100,
			Y:      100,
			Width:  250,
			Height: 150,
			Label:  label,
			Shape:  "rectangle",
		},
		Style: schema.Style{
			FillColor:     "#E3F2FD",
			StrokeColor:   "#1976D2",
			StrokeWidth:   2,
			Rounded:       true,
			FontSize:      12,
			TextAlign:     "center",
			VerticalAlign: "middle",
		},
	}, nil
}

// getStringParam safely gets a string parameter with a default value
func (p *AWSProvider) getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if value, ok := params[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}
