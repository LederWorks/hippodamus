package providers

// InitializeBuiltinProviders registers all built-in providers with the given version
// This function should be called from main.go to avoid import cycles
func InitializeBuiltinProviders(version string) error {
	// This function is intentionally empty - providers should be registered
	// from main.go to avoid import cycle issues.
	// See main.go for the actual provider registration code.
	return nil
}
