package aws

import (
	"testing"

	"github.com/LederWorks/hippodamus/pkg/providers"
	"github.com/LederWorks/hippodamus/pkg/schema"
)

func TestAWSProvider_Basic(t *testing.T) {
	provider := NewAWSProvider()

	if provider.Name() != "aws" {
		t.Errorf("Expected provider name 'aws', got '%s'", provider.Name())
	}

	if provider.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", provider.Version())
	}
}

func TestAWSProvider_Resources(t *testing.T) {
	provider := NewAWSProvider()
	resources := provider.Resources()

	if len(resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resources))
	}

	// Check organization resource
	orgResource := resources[0]
	if orgResource.Type != "organization" {
		t.Errorf("Expected first resource type 'organization', got '%s'", orgResource.Type)
	}

	if orgResource.Name != "AWS Organization" {
		t.Errorf("Expected first resource name 'AWS Organization', got '%s'", orgResource.Name)
	}

	// Check VPC resource
	vpcResource := resources[1]
	if vpcResource.Type != "vpc" {
		t.Errorf("Expected second resource type 'vpc', got '%s'", vpcResource.Type)
	}
}

func TestAWSProvider_ValidateOrganization(t *testing.T) {
	provider := NewAWSProvider()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid organization",
			params: map[string]interface{}{
				"orgName":             "Test Org",
				"managementAccountId": "123456789012",
			},
			wantErr: false,
		},
		{
			name: "missing orgName",
			params: map[string]interface{}{
				"managementAccountId": "123456789012",
			},
			wantErr: true,
		},
		{
			name: "missing managementAccountId",
			params: map[string]interface{}{
				"orgName": "Test Org",
			},
			wantErr: true,
		},
		{
			name:    "empty params",
			params:  map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.Validate("organization", tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAWSProvider_ValidateVPC(t *testing.T) {
	provider := NewAWSProvider()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid VPC",
			params: map[string]interface{}{
				"vpcName":   "test-vpc",
				"cidrBlock": "10.0.0.0/16",
			},
			wantErr: false,
		},
		{
			name: "missing vpcName",
			params: map[string]interface{}{
				"cidrBlock": "10.0.0.0/16",
			},
			wantErr: true,
		},
		{
			name: "missing cidrBlock",
			params: map[string]interface{}{
				"vpcName": "test-vpc",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.Validate("vpc", tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAWSProvider_GenerateOrganizationTemplate(t *testing.T) {
	provider := NewAWSProvider()

	params := map[string]interface{}{
		"orgName":             "Test Organization",
		"managementAccountId": "123456789012",
		"fillColor":           "#FFF8E1",
		"strokeColor":         "#FF9900",
	}

	element, err := provider.GenerateTemplate("organization", params)
	if err != nil {
		t.Errorf("GenerateTemplate() error = %v", err)
		return
	}

	if element.Type != schema.ElementTypeShape {
		t.Errorf("Expected element type %s, got %s", schema.ElementTypeShape, element.Type)
	}

	if element.Properties.Label != "Test Organization" {
		t.Errorf("Expected label 'Test Organization', got '%s'", element.Properties.Label)
	}

	if element.Style.FillColor != "#FFF8E1" {
		t.Errorf("Expected fill color '#FFF8E1', got '%s'", element.Style.FillColor)
	}

	if element.Nesting.Mode != schema.NestingModeChild {
		t.Errorf("Expected nesting mode %s, got %s", schema.NestingModeChild, element.Nesting.Mode)
	}
}

func TestAWSProvider_GenerateVPCTemplate(t *testing.T) {
	provider := NewAWSProvider()

	params := map[string]interface{}{
		"vpcName":   "production-vpc",
		"cidrBlock": "10.0.0.0/16",
		"region":    "us-west-2",
	}

	element, err := provider.GenerateTemplate("vpc", params)
	if err != nil {
		t.Errorf("GenerateTemplate() error = %v", err)
		return
	}

	if element.Type != schema.ElementTypeShape {
		t.Errorf("Expected element type %s, got %s", schema.ElementTypeShape, element.Type)
	}

	expectedLabel := "production-vpc\\n10.0.0.0/16\\nus-west-2"
	if element.Properties.Label != expectedLabel {
		t.Errorf("Expected label '%s', got '%s'", expectedLabel, element.Properties.Label)
	}
}

func TestAWSProvider_UnsupportedResource(t *testing.T) {
	provider := NewAWSProvider()

	err := provider.Validate("unsupported", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error for unsupported resource type")
	}

	if providerErr, ok := err.(*providers.ProviderError); ok {
		if providerErr.Code != "UNSUPPORTED_RESOURCE" {
			t.Errorf("Expected error code 'UNSUPPORTED_RESOURCE', got '%s'", providerErr.Code)
		}
	} else {
		t.Error("Expected ProviderError type")
	}
}

func TestAWSProvider_GetSchema(t *testing.T) {
	provider := NewAWSProvider()

	schema, err := provider.GetSchema("organization")
	if err != nil {
		t.Errorf("GetSchema() error = %v", err)
		return
	}

	if schema == nil {
		t.Error("Expected schema to be non-nil")
		return
	}

	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	// Test non-existent resource
	_, err = provider.GetSchema("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent resource schema")
	}
}
