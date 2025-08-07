package providers

// GetBuiltinProviders returns a list of all built-in provider names
func GetBuiltinProviders() []string {
	return []string{
		"aws",
		// TODO: Add "azure", "gcp" when implemented
	}
}
