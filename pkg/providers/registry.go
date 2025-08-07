package providers

import (
	"fmt"
	"sort"
	"sync"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

// Registry manages all available providers
type Registry struct {
	providers map[string]Provider
	mutex     sync.RWMutex
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register registers a provider with the registry
func (r *Registry) Register(provider Provider) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := provider.Name()
	if name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %s is already registered", name)
	}

	r.providers[name] = provider
	return nil
}

// Get retrieves a provider by name
func (r *Registry) Get(name string) (Provider, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// List returns all registered provider names
func (r *Registry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// GetAll returns all registered providers
func (r *Registry) GetAll() map[string]Provider {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	providers := make(map[string]Provider)
	for name, provider := range r.providers {
		providers[name] = provider
	}

	return providers
}

// Unregister removes a provider from the registry
func (r *Registry) Unregister(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.providers[name]; !exists {
		return fmt.Errorf("provider %s not found", name)
	}

	delete(r.providers, name)
	return nil
}

// GetResourceTypes returns all resource types across all providers
func (r *Registry) GetResourceTypes() map[string]ResourceDefinition {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	resources := make(map[string]ResourceDefinition)
	for _, provider := range r.providers {
		for _, resource := range provider.Resources() {
			// Prefix resource type with provider name to avoid conflicts
			key := fmt.Sprintf("%s-%s", provider.Name(), resource.Type)
			resources[key] = resource
		}
	}

	return resources
}

// GetResourcesByProvider returns all resources for a specific provider
func (r *Registry) GetResourcesByProvider(providerName string) ([]ResourceDefinition, error) {
	provider, err := r.Get(providerName)
	if err != nil {
		return nil, err
	}

	return provider.Resources(), nil
}

// ValidateResource validates a resource configuration using the appropriate provider
func (r *Registry) ValidateResource(providerName, resourceType string, params map[string]interface{}) error {
	provider, err := r.Get(providerName)
	if err != nil {
		return err
	}

	return provider.Validate(resourceType, params)
}

// GenerateTemplate generates a template using the appropriate provider
func (r *Registry) GenerateTemplate(providerName, resourceType string, params map[string]interface{}) (*schema.Element, error) {
	provider, err := r.Get(providerName)
	if err != nil {
		return nil, err
	}

	return provider.GenerateTemplate(resourceType, params)
}

// DefaultRegistry is the global provider registry instance
var DefaultRegistry = NewRegistry()
