package schema

import (
	"testing"
	"time"
)

func TestDiagramConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  DiagramConfig
		wantErr bool
	}{
		{
			name: "valid minimal config",
			config: DiagramConfig{
				Version: "1.0",
				Diagram: Diagram{
					Pages: []Page{
						{
							ID:   "page1",
							Name: "Test Page",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			config: DiagramConfig{
				Diagram: Diagram{
					Pages: []Page{
						{
							ID:   "page1",
							Name: "Test Page",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing pages",
			config: DiagramConfig{
				Version: "1.0",
				Diagram: Diagram{
					Pages: []Page{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDiagramConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDiagramConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetadata_Defaults(t *testing.T) {
	metadata := Metadata{
		Title:   "Test Diagram",
		Created: time.Now(),
	}

	if metadata.Title != "Test Diagram" {
		t.Errorf("Expected title 'Test Diagram', got '%s'", metadata.Title)
	}

	if metadata.Created.IsZero() {
		t.Error("Expected created time to be set")
	}
}

func TestElementType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		elemType ElementType
		expected string
	}{
		{"shape", ElementTypeShape, "shape"},
		{"connector", ElementTypeConnector, "connector"},
		{"text", ElementTypeText, "text"},
		{"group", ElementTypeGroup, "group"},
		{"swimlane", ElementTypeSwimLane, "swimlane"},
		{"template", ElementTypeTemplate, "template"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.elemType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.elemType))
			}
		})
	}
}

// validateDiagramConfig validates a diagram configuration
func validateDiagramConfig(config *DiagramConfig) error {
	if config.Version == "" {
		return &ValidationError{Field: "version", Message: "version is required"}
	}

	if len(config.Diagram.Pages) == 0 {
		return &ValidationError{Field: "diagram.pages", Message: "at least one page is required"}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
