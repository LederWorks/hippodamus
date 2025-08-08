package providers

import (
	"fmt"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

// Provider defines the interface that all Hippodamus providers must implement
type Provider interface {
	// Name returns the provider name (e.g., "aws", "azure", "gcp")
	Name() string

	// Version returns the provider version (semantic versioning)
	Version() string

	// Resources returns the list of resources this provider supports
	Resources() []ResourceDefinition

	// GenerateTemplate generates a template element for the given resource type
	GenerateTemplate(resourceType string, params map[string]interface{}) (*schema.Element, error)

	// Validate validates the provider configuration and parameters
	Validate(resourceType string, params map[string]interface{}) error

	// GetSchema returns the JSON schema for a specific resource type
	GetSchema(resourceType string) (map[string]interface{}, error)
}

// ResourceDefinition defines a resource type that a provider supports
type ResourceDefinition struct {
	Type        string                 `json:"type"`        // Resource type (e.g., "aws-vpc", "azure-rg")
	Name        string                 `json:"name"`        // Human-readable name
	Description string                 `json:"description"` // Resource description
	Category    string                 `json:"category"`    // Resource category (compute, network, storage, etc.)
	Schema      map[string]interface{} `json:"schema"`      // JSON schema for validation
	Examples    []ResourceExample      `json:"examples"`    // Usage examples
}

// ResourceExample provides usage examples for a resource
type ResourceExample struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

// ProviderConfig contains provider initialization configuration
type ProviderConfig struct {
	Name     string                 `json:"name"`
	Version  string                 `json:"version"`
	Settings map[string]interface{} `json:"settings"`
}

// ProviderMetadata contains metadata about a provider
type ProviderMetadata struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Author       string   `json:"author"`
	Homepage     string   `json:"homepage"`
	Repository   string   `json:"repository"`
	Dependencies []string `json:"dependencies"`
	Tags         []string `json:"tags"`
}

// ValidationError represents a provider validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ProviderError represents a provider-specific error
type ProviderError struct {
	Provider string `json:"provider"`
	Resource string `json:"resource"`
	Message  string `json:"message"`
	Code     string `json:"code"`
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider %s: %s", e.Provider, e.Message)
}
