package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/LederWorks/hippodamus/pkg/schema"
)

// TemplateProcessor handles template loading and processing with hive support
type TemplateProcessor struct {
	templates   map[string]*schema.Template
	templateDir string
	hives       map[string][]string // Maps hive name to list of templates
}

// NewTemplateProcessor creates a new template processor with hive support
func NewTemplateProcessor(templateDir string) *TemplateProcessor {
	return &TemplateProcessor{
		templates:   make(map[string]*schema.Template),
		templateDir: templateDir,
		hives:       make(map[string][]string),
	}
}

// getElementDisplayName returns the display name for an element (used for error messages)
func (tp *TemplateProcessor) getElementDisplayName(element *schema.Element) string {
	if element.Name != "" {
		return element.Name
	}
	return element.ID
}

// LoadTemplates loads all templates from the configured template directory with hive support
func (tp *TemplateProcessor) LoadTemplates() error {
	if tp.templateDir == "" {
		return nil // No templates directory specified
	}

	return filepath.Walk(tp.templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			template, err := tp.loadTemplate(path)
			if err != nil {
				return fmt.Errorf("failed to load template %s: %w", path, err)
			}

			// Determine hive from file path
			hive := tp.getHiveFromPath(path)
			templateKey := tp.getTemplateKey(template.Name, hive)

			tp.templates[templateKey] = template

			// Register template in hive
			if hive != "" {
				if tp.hives[hive] == nil {
					tp.hives[hive] = make([]string, 0)
				}
				tp.hives[hive] = append(tp.hives[hive], template.Name)
			}
		}

		return nil
	})
}

// LoadTemplateRefs loads templates from template references
func (tp *TemplateProcessor) LoadTemplateRefs(refs []schema.TemplateRef) error {
	for _, ref := range refs {
		templatePath := ref.Template

		// Handle relative paths
		if !filepath.IsAbs(templatePath) {
			if tp.templateDir != "" {
				// If template path already contains "templates/", don't double it
				if strings.HasPrefix(templatePath, "templates/") || strings.HasPrefix(templatePath, "templates\\") {
					templatePath = filepath.Join(tp.templateDir, "..", templatePath)
				} else {
					templatePath = filepath.Join(tp.templateDir, templatePath)
				}
			}
		}

		template, err := tp.loadTemplate(templatePath)
		if err != nil {
			return fmt.Errorf("failed to load template %s from %s: %w", ref.Name, templatePath, err)
		}

		// Override the template name with the reference name if different
		if ref.Name != template.Name {
			template.Name = ref.Name
		}

		tp.templates[template.Name] = template
	}

	return nil
}

// loadTemplate loads a single template from a file
func (tp *TemplateProcessor) loadTemplate(path string) (*schema.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var template schema.Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, err
	}

	return &template, nil
}

// ProcessDiagram processes a diagram configuration and applies templates
func (tp *TemplateProcessor) ProcessDiagram(config *schema.DiagramConfig) error {
	// Load template references first
	if err := tp.LoadTemplateRefs(config.Templates); err != nil {
		return err
	}

	// Process each page
	for i := range config.Diagram.Pages {
		if err := tp.processPage(&config.Diagram.Pages[i]); err != nil {
			return fmt.Errorf("failed to process page %s: %w", config.Diagram.Pages[i].ID, err)
		}
	}

	return nil
}

// processPage processes a page and applies templates to its elements
func (tp *TemplateProcessor) processPage(page *schema.Page) error {
	// Process layers
	for i := range page.Layers {
		if err := tp.processElements(page.Layers[i].Elements); err != nil {
			return fmt.Errorf("failed to process layer %s: %w", page.Layers[i].ID, err)
		}
	}

	// Process page-level elements
	if err := tp.processElements(page.Elements); err != nil {
		return fmt.Errorf("failed to process page elements: %w", err)
	}

	return nil
}

// processElements processes a list of elements and applies templates
func (tp *TemplateProcessor) processElements(elements []schema.Element) error {
	return tp.processElementsWithContext(elements, []string{})
}

// processElementsWithContext processes elements with parent template context
func (tp *TemplateProcessor) processElementsWithContext(elements []schema.Element, parentTemplates []string) error {
	for i := range elements {
		if err := tp.processElementWithContext(&elements[i], parentTemplates); err != nil {
			return fmt.Errorf("failed to process element %s: %w", elements[i].ID, err)
		}

		// Process children recursively with updated parent context
		var childContext []string
		if elements[i].Template != "" {
			childContext = append([]string{elements[i].Template}, parentTemplates...)
		} else {
			childContext = parentTemplates
		}

		if err := tp.processElementsWithContext(elements[i].Children, childContext); err != nil {
			return fmt.Errorf("failed to process children of element %s: %w", elements[i].ID, err)
		}
	}

	return nil
}

// processElement processes a single element and applies its template if specified
func (tp *TemplateProcessor) processElement(element *schema.Element) error {
	return tp.processElementWithContext(element, []string{})
}

// processElementWithContext processes a single element with parent template context
func (tp *TemplateProcessor) processElementWithContext(element *schema.Element, parentTemplates []string) error {
	// Handle template type - inherit type from template
	if element.Type == schema.ElementTypeTemplate {
		if element.Template == "" {
			return fmt.Errorf("element %s has type 'template' but no template specified", tp.getElementDisplayName(element))
		}

		// Resolve template reference using hive-aware resolution
		currentHive := tp.getCurrentHive(parentTemplates)
		resolvedTemplate := tp.resolveTemplateReference(element.Template, currentHive)

		template, exists := tp.templates[resolvedTemplate]
		if !exists {
			return fmt.Errorf("template %s not found for element %s (resolved to: %s)", element.Template, tp.getElementDisplayName(element), resolvedTemplate)
		}

		// Validate template dependencies
		if err := tp.validateDependencies(element, template, parentTemplates); err != nil {
			return fmt.Errorf("dependency validation failed for element %s: %w", tp.getElementDisplayName(element), err)
		}

		// Inherit the type from the template's first element
		if len(template.Elements) > 0 {
			element.Type = template.Elements[0].Type
		} else {
			return fmt.Errorf("template %s has no elements to inherit type from", element.Template)
		}
	}

	if element.Template == "" {
		return nil // No template to apply
	}

	// Resolve template reference for non-template type elements too
	currentHive := tp.getCurrentHive(parentTemplates)
	resolvedTemplate := tp.resolveTemplateReference(element.Template, currentHive)

	template, exists := tp.templates[resolvedTemplate]
	if !exists {
		return fmt.Errorf("template %s not found (resolved to: %s)", element.Template, resolvedTemplate)
	}

	// Apply template to element
	if err := tp.applyTemplate(element, template); err != nil {
		return fmt.Errorf("failed to apply template %s to element %s: %w", element.Template, tp.getElementDisplayName(element), err)
	}

	return nil
}

// applyTemplate applies a template to an element
func (tp *TemplateProcessor) applyTemplate(element *schema.Element, tmpl *schema.Template) error {
	// Prepare template variables
	vars := make(map[string]interface{})

	// Add element properties as variables
	vars["id"] = element.ID
	vars["name"] = element.Name
	vars["x"] = element.Properties.X
	vars["y"] = element.Properties.Y
	vars["width"] = element.Properties.Width
	vars["height"] = element.Properties.Height
	vars["label"] = element.Properties.Label

	// Add custom properties
	for key, value := range element.Properties.Custom {
		vars[key] = value
	}

	// Set default values for template parameters
	for _, param := range tmpl.Parameters {
		if _, exists := vars[param.Name]; !exists && param.Default != nil {
			vars[param.Name] = param.Default
		}
	}

	// Add specific default values for common template variables
	if _, exists := vars["fillColor"]; !exists {
		vars["fillColor"] = "#E3F2FD"
	}
	if _, exists := vars["strokeColor"]; !exists {
		vars["strokeColor"] = "#1976D2"
	}

	// Validate required parameters
	for _, param := range tmpl.Parameters {
		if param.Required {
			if _, exists := vars[param.Name]; !exists {
				return fmt.Errorf("required template parameter %s not provided", param.Name)
			}
		}
	}

	// Apply template elements
	if len(tmpl.Elements) > 0 {
		// If template has multiple elements, create them as children
		if len(tmpl.Elements) > 1 {
			element.Children = make([]schema.Element, 0, len(tmpl.Elements))
		}

		for i, templateElement := range tmpl.Elements {
			var targetElement *schema.Element

			if i == 0 && len(tmpl.Elements) == 1 {
				// Single template element: merge with current element
				targetElement = element
			} else {
				// Multiple template elements: create as children
				childElement := templateElement
				childElement.ID = fmt.Sprintf("%s-%d", element.ID, i)
				element.Children = append(element.Children, childElement)
				targetElement = &element.Children[len(element.Children)-1]
			}

			// Merge template element properties first
			tp.mergeElementProperties(targetElement, &templateElement)

			// Then apply template variables to element
			if err := tp.applyTemplateVariables(targetElement, vars); err != nil {
				return fmt.Errorf("failed to apply template variables: %w", err)
			}

			// Process children that were merged from template for this specific element
			if len(targetElement.Children) > 0 {
				// Create a copy of vars to avoid modification issues
				childVars := make(map[string]interface{})
				for k, v := range vars {
					childVars[k] = v
				}
				for i := range targetElement.Children {
					if err := tp.applyTemplateVariables(&targetElement.Children[i], childVars); err != nil {
						return fmt.Errorf("failed to apply template variables to child %d: %w", i, err)
					}
				}
			}
		}
	}

	return nil
}

// applyTemplateVariables applies template variables to an element using Go's template engine
func (tp *TemplateProcessor) applyTemplateVariables(element *schema.Element, vars map[string]interface{}) error {
	// Apply variables to label
	if element.Properties.Label != "" {
		label, err := tp.processTemplateString(element.Properties.Label, vars)
		if err != nil {
			return fmt.Errorf("failed to process label template: %w", err)
		}
		element.Properties.Label = label
	}

	// Apply variables to value
	if element.Properties.Value != "" {
		value, err := tp.processTemplateString(element.Properties.Value, vars)
		if err != nil {
			return fmt.Errorf("failed to process value template: %w", err)
		}
		element.Properties.Value = value
	}

	// Apply variables to shape
	if element.Properties.Shape != "" {
		shape, err := tp.processTemplateString(element.Properties.Shape, vars)
		if err != nil {
			return fmt.Errorf("failed to process shape template: %w", err)
		}
		element.Properties.Shape = shape
	}

	// Apply template variables to style properties
	if element.Style.FillColor != "" {
		fillColor, err := tp.processTemplateString(element.Style.FillColor, vars)
		if err != nil {
			return fmt.Errorf("failed to process fillColor template: %w", err)
		}
		element.Style.FillColor = fillColor
	}

	if element.Style.StrokeColor != "" {
		strokeColor, err := tp.processTemplateString(element.Style.StrokeColor, vars)
		if err != nil {
			return fmt.Errorf("failed to process strokeColor template: %w", err)
		}
		element.Style.StrokeColor = strokeColor
	}

	if element.Style.FontFamily != "" {
		fontFamily, err := tp.processTemplateString(element.Style.FontFamily, vars)
		if err != nil {
			return fmt.Errorf("failed to process fontFamily template: %w", err)
		}
		element.Style.FontFamily = fontFamily
	}

	if element.Style.FontColor != "" {
		fontColor, err := tp.processTemplateString(element.Style.FontColor, vars)
		if err != nil {
			return fmt.Errorf("failed to process fontColor template: %w", err)
		}
		element.Style.FontColor = fontColor
	}

	if element.Style.FontStyle != "" {
		fontStyle, err := tp.processTemplateString(element.Style.FontStyle, vars)
		if err != nil {
			return fmt.Errorf("failed to process fontStyle template: %w", err)
		}
		element.Style.FontStyle = fontStyle
	}

	if element.Style.TextAlign != "" {
		textAlign, err := tp.processTemplateString(element.Style.TextAlign, vars)
		if err != nil {
			return fmt.Errorf("failed to process textAlign template: %w", err)
		}
		element.Style.TextAlign = textAlign
	}

	if element.Style.VerticalAlign != "" {
		verticalAlign, err := tp.processTemplateString(element.Style.VerticalAlign, vars)
		if err != nil {
			return fmt.Errorf("failed to process verticalAlign template: %w", err)
		}
		element.Style.VerticalAlign = verticalAlign
	}

	if element.Style.LabelPosition != "" {
		labelPosition, err := tp.processTemplateString(element.Style.LabelPosition, vars)
		if err != nil {
			return fmt.Errorf("failed to process labelPosition template: %w", err)
		}
		element.Style.LabelPosition = labelPosition
	}

	if element.Style.VerticalLabelPosition != "" {
		verticalLabelPosition, err := tp.processTemplateString(element.Style.VerticalLabelPosition, vars)
		if err != nil {
			return fmt.Errorf("failed to process verticalLabelPosition template: %w", err)
		}
		element.Style.VerticalLabelPosition = verticalLabelPosition
	}

	if element.Style.StrokeDashArray != "" {
		strokeDashArray, err := tp.processTemplateString(element.Style.StrokeDashArray, vars)
		if err != nil {
			return fmt.Errorf("failed to process strokeDashArray template: %w", err)
		}
		element.Style.StrokeDashArray = strokeDashArray
	}

	// Apply variables to custom style properties
	for key, value := range element.Style.Custom {
		processedValue, err := tp.processTemplateString(value, vars)
		if err != nil {
			return fmt.Errorf("failed to process custom style %s: %w", key, err)
		}
		element.Style.Custom[key] = processedValue
	}

	return nil
}

// processTemplateString processes a template string with variables
func (tp *TemplateProcessor) processTemplateString(templateStr string, vars map[string]interface{}) (string, error) {
	// Define template functions
	funcMap := template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			return a != b
		},
		"and": func(a, b bool) bool {
			return a && b
		},
		"or": func(a, b bool) bool {
			return a || b
		},
		"not": func(a bool) bool {
			return !a
		},
	}

	tmpl, err := template.New("template").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, vars); err != nil {
		return "", err
	}

	return result.String(), nil
}

// mergeElementProperties merges properties from template element to target element
func (tp *TemplateProcessor) mergeElementProperties(target *schema.Element, template *schema.Element) {
	// Merge type if not set
	if target.Type == "" {
		target.Type = template.Type
	}

	// Merge properties
	if target.Properties.Width == 0 && template.Properties.Width > 0 {
		target.Properties.Width = template.Properties.Width
	}
	if target.Properties.Height == 0 && template.Properties.Height > 0 {
		target.Properties.Height = template.Properties.Height
	}
	if target.Properties.Shape == "" {
		target.Properties.Shape = template.Properties.Shape
	}
	if target.Properties.ShapeType == "" {
		target.Properties.ShapeType = template.Properties.ShapeType
	}

	// Merge style properties
	if target.Style.FillColor == "" {
		target.Style.FillColor = template.Style.FillColor
	}
	if target.Style.StrokeColor == "" {
		target.Style.StrokeColor = template.Style.StrokeColor
	}
	if target.Style.StrokeWidth == 0 && template.Style.StrokeWidth > 0 {
		target.Style.StrokeWidth = template.Style.StrokeWidth
	}
	if target.Style.StrokeDashArray == "" {
		target.Style.StrokeDashArray = template.Style.StrokeDashArray
	}
	if target.Style.FontFamily == "" {
		target.Style.FontFamily = template.Style.FontFamily
	}
	if target.Style.FontSize == 0 && template.Style.FontSize > 0 {
		target.Style.FontSize = template.Style.FontSize
	}
	if target.Style.FontColor == "" {
		target.Style.FontColor = template.Style.FontColor
	}
	if target.Style.FontStyle == "" {
		target.Style.FontStyle = template.Style.FontStyle
	}
	if target.Style.TextAlign == "" {
		target.Style.TextAlign = template.Style.TextAlign
	}
	if target.Style.VerticalAlign == "" {
		target.Style.VerticalAlign = template.Style.VerticalAlign
	}
	if target.Style.LabelPosition == "" {
		target.Style.LabelPosition = template.Style.LabelPosition
	}
	if target.Style.VerticalLabelPosition == "" {
		target.Style.VerticalLabelPosition = template.Style.VerticalLabelPosition
	}
	if !target.Style.Rounded && template.Style.Rounded {
		target.Style.Rounded = template.Style.Rounded
	}
	if !target.Style.Shadow && template.Style.Shadow {
		target.Style.Shadow = template.Style.Shadow
	}
	if !target.Style.Glass && template.Style.Glass {
		target.Style.Glass = template.Style.Glass
	}

	// Merge custom style properties
	if target.Style.Custom == nil {
		target.Style.Custom = make(map[string]string)
	}
	for key, value := range template.Style.Custom {
		if _, exists := target.Style.Custom[key]; !exists {
			target.Style.Custom[key] = value
		}
	}

	// Merge tags
	if len(target.Tags) == 0 {
		target.Tags = template.Tags
	}

	// Merge children if target has none and template has children
	if len(target.Children) == 0 && len(template.Children) > 0 {
		target.Children = make([]schema.Element, len(template.Children))
		// Deep copy children to avoid sharing references
		for i, child := range template.Children {
			target.Children[i] = child // This creates a copy of the struct
		}
	}

	// Merge nesting configuration
	if target.Nesting.Mode == "" {
		target.Nesting = template.Nesting
	}
}

// GetTemplate returns a template by name
func (tp *TemplateProcessor) GetTemplate(name string) (*schema.Template, bool) {
	template, exists := tp.templates[name]
	return template, exists
}

// ListTemplates returns all loaded template names
func (tp *TemplateProcessor) ListTemplates() []string {
	names := make([]string, 0, len(tp.templates))
	for name := range tp.templates {
		names = append(names, name)
	}
	return names
}

// validateDependencies validates that an element's template dependencies are satisfied
func (tp *TemplateProcessor) validateDependencies(element *schema.Element, template *schema.Template, parentTemplates []string) error {
	// If template has no dependencies, validation passes
	if len(template.Dependencies) == 0 {
		return nil
	}

	// Check each required dependency
	for _, dep := range template.Dependencies {
		if !dep.Required {
			continue // Skip optional dependencies
		}

		switch dep.Relationship {
		case "parent":
			if !tp.hasParentOfType(parentTemplates, dep.Type) {
				return fmt.Errorf("required parent dependency not satisfied: template %s requires a parent of type %s, but element %s has parents: %v",
					template.Name, dep.Type, tp.getElementDisplayName(element), parentTemplates)
			}
		case "ancestor":
			if !tp.hasAncestorOfType(parentTemplates, dep.Type) {
				return fmt.Errorf("required ancestor dependency not satisfied: template %s requires an ancestor of type %s, but element %s has ancestors: %v",
					template.Name, dep.Type, tp.getElementDisplayName(element), parentTemplates)
			}
		}
	}

	return nil
}

// hasParentOfType checks if there's a direct parent of the specified template type
func (tp *TemplateProcessor) hasParentOfType(parentTemplates []string, templateType string) bool {
	if len(parentTemplates) == 0 {
		return false
	}
	// Direct parent is the first in the list
	return parentTemplates[0] == templateType
}

// hasAncestorOfType checks if there's any ancestor of the specified template type
func (tp *TemplateProcessor) hasAncestorOfType(parentTemplates []string, templateType string) bool {
	for _, parent := range parentTemplates {
		if parent == templateType {
			return true
		}
	}
	return false
}

// getCurrentHive determines the current hive context from parent templates
func (tp *TemplateProcessor) getCurrentHive(parentTemplates []string) string {
	if len(parentTemplates) == 0 {
		return ""
	}

	// Get the hive from the most recent (first) parent template
	parentTemplate := parentTemplates[0]

	// Check if parent template contains hive notation
	if strings.Contains(parentTemplate, "/") {
		parts := strings.Split(parentTemplate, "/")
		return parts[0]
	}

	return ""
}

// getHiveFromPath extracts the hive name from a template file path
func (tp *TemplateProcessor) getHiveFromPath(templatePath string) string {
	// Convert to forward slashes for consistent handling
	templatePath = filepath.ToSlash(templatePath)
	templateDir := filepath.ToSlash(tp.templateDir)

	// Get relative path from template directory
	relPath, err := filepath.Rel(templateDir, templatePath)
	if err != nil {
		return "" // Not in template directory, no hive
	}

	// Convert back to forward slashes
	relPath = filepath.ToSlash(relPath)

	// Check if template is in a subdirectory (hive)
	parts := strings.Split(relPath, "/")
	if len(parts) > 1 {
		return parts[0] // First directory is the hive
	}

	return "" // Template is in root, no hive
}

// getTemplateKey generates a unique key for template storage
func (tp *TemplateProcessor) getTemplateKey(templateName, hive string) string {
	if hive == "" {
		return templateName
	}
	return hive + "/" + templateName
}

// resolveTemplateReference resolves a template reference, supporting hive notation
func (tp *TemplateProcessor) resolveTemplateReference(templateRef string, currentHive string) string {
	// If template reference contains a slash, it's already fully qualified
	if strings.Contains(templateRef, "/") {
		return templateRef
	}

	// Try current hive first if we're in one
	if currentHive != "" {
		hiveKey := tp.getTemplateKey(templateRef, currentHive)
		if _, exists := tp.templates[hiveKey]; exists {
			return hiveKey
		}
	}

	// Try root level
	if _, exists := tp.templates[templateRef]; exists {
		return templateRef
	}

	// Return original reference for error handling
	return templateRef
}

// ListHives returns all available template hives
func (tp *TemplateProcessor) ListHives() []string {
	hives := make([]string, 0, len(tp.hives))
	for hive := range tp.hives {
		hives = append(hives, hive)
	}
	return hives
}

// ListTemplatesInHive returns all templates in a specific hive
func (tp *TemplateProcessor) ListTemplatesInHive(hive string) []string {
	if templates, exists := tp.hives[hive]; exists {
		return templates
	}
	return []string{}
}
