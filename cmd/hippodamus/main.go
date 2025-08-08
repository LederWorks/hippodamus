package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/LederWorks/hippodamus/pkg/drawio"
	"github.com/LederWorks/hippodamus/pkg/providers"
	"github.com/LederWorks/hippodamus/pkg/schema"
	"github.com/LederWorks/hippodamus/pkg/templates"
	"github.com/LederWorks/hippodamus/providers/core"
)

// Version information injected at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

type Config struct {
	InputFile     string
	OutputFile    string
	TemplatesDir  string
	ValidateOnly  bool
	ShowVersion   bool
	ListProviders bool
	Verbose       bool
}

func main() {
	config := parseFlags()

	// Initialize built-in providers with the current application version
	if err := initializeProviders(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing providers: %v\n", err)
		os.Exit(1)
	}

	if config.ShowVersion {
		fmt.Printf("Hippodamus v%s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
		return
	}

	if config.ListProviders {
		listProviders()
		return
	}

	if config.InputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: input file is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.InputFile, "input", "", "Input YAML file path")
	flag.StringVar(&config.InputFile, "i", "", "Input YAML file path (short form)")
	flag.StringVar(&config.OutputFile, "output", "", "Output file path (.xml or .drawio, default: input file with .drawio extension)")
	flag.StringVar(&config.OutputFile, "o", "", "Output file path (short form)")
	flag.StringVar(&config.TemplatesDir, "templates", "", "Templates directory path")
	flag.StringVar(&config.TemplatesDir, "t", "", "Templates directory path (short form)")
	flag.BoolVar(&config.ValidateOnly, "validate", false, "Validate YAML only, don't generate output")
	flag.BoolVar(&config.ValidateOnly, "v", false, "Validate YAML only (short form)")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version information")
	flag.BoolVar(&config.ListProviders, "list-providers", false, "List available providers and their resources")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Hippodamus v%s - YAML to Draw.io XML Converter\n\n", version)
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSupported output formats:\n")
		fmt.Fprintf(os.Stderr, "  .xml     - Standard XML format\n")
		fmt.Fprintf(os.Stderr, "  .drawio  - Draw.io native format\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -i diagram.yaml -o diagram.xml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i diagram.yaml -o diagram.drawio\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -input diagram.yaml -templates ./templates\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -validate -input diagram.yaml\n", os.Args[0])
	}

	flag.Parse()

	// Set default output file if not provided
	if config.OutputFile == "" && config.InputFile != "" && !config.ValidateOnly {
		ext := filepath.Ext(config.InputFile)
		config.OutputFile = config.InputFile[:len(config.InputFile)-len(ext)] + ".drawio"
	}

	return config
}

func run(config *Config) error {
	if config.Verbose {
		fmt.Printf("Loading YAML configuration from: %s\n", config.InputFile)
	}

	// Load and parse YAML configuration
	diagramConfig, err := loadDiagramConfig(config.InputFile)
	if err != nil {
		return fmt.Errorf("failed to load diagram configuration: %w", err)
	}

	if config.Verbose {
		fmt.Printf("Loaded diagram: %s (version %s)\n", diagramConfig.Metadata.Title, diagramConfig.Version)
		fmt.Printf("Pages: %d\n", len(diagramConfig.Diagram.Pages))
	}

	// Initialize template processor
	templateProcessor := templates.NewTemplateProcessor(config.TemplatesDir)

	// Load templates if template directory is specified
	if config.TemplatesDir != "" {
		if config.Verbose {
			fmt.Printf("Loading templates from: %s\n", config.TemplatesDir)
		}
		if err := templateProcessor.LoadTemplates(); err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}

		if config.Verbose {
			templateNames := templateProcessor.ListTemplates()
			fmt.Printf("Loaded %d templates: %v\n", len(templateNames), templateNames)
		}
	}

	// Process diagram with templates
	if err := templateProcessor.ProcessDiagram(diagramConfig); err != nil {
		return fmt.Errorf("failed to process diagram templates: %w", err)
	}

	if config.ValidateOnly {
		fmt.Println("YAML configuration is valid")
		return nil
	}

	if config.Verbose {
		fmt.Printf("Generating draw.io XML output\n")
	}

	// Generate draw.io XML
	generator := drawio.NewGenerator()
	document, err := generator.Generate(diagramConfig)
	if err != nil {
		return fmt.Errorf("failed to generate draw.io XML: %w", err)
	}

	if config.Verbose {
		fmt.Printf("Writing output to: %s\n", config.OutputFile)
	}

	// Write output file
	if err := writeDrawioXML(document, config.OutputFile); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", config.InputFile, config.OutputFile)
	return nil
}

func loadDiagramConfig(filename string) (*schema.DiagramConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config schema.DiagramConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Validate required fields
	if config.Version == "" {
		return nil, fmt.Errorf("version field is required")
	}

	if len(config.Diagram.Pages) == 0 {
		return nil, fmt.Errorf("at least one page is required")
	}

	// Validate page IDs are unique
	pageIDs := make(map[string]bool)
	for _, page := range config.Diagram.Pages {
		if page.ID == "" {
			return nil, fmt.Errorf("page ID is required")
		}
		if pageIDs[page.ID] {
			return nil, fmt.Errorf("duplicate page ID: %s", page.ID)
		}
		pageIDs[page.ID] = true
	}

	return &config, nil
}

func writeDrawioXML(document *drawio.DrawioDocument, filename string) error {
	// Create output directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write XML declaration
	if _, err := file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"); err != nil {
		return err
	}

	// Marshal and write XML
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	if err := encoder.Encode(document); err != nil {
		return err
	}

	return nil
}

// listProviders displays all available providers and their resources
func listProviders() {
	fmt.Println("🔧 Available Providers")
	fmt.Println("======================")

	registry := providers.DefaultRegistry
	providerNames := registry.List()

	if len(providerNames) == 0 {
		fmt.Println("No providers registered.")
		return
	}

	for _, name := range providerNames {
		provider, err := registry.Get(name)
		if err != nil {
			fmt.Printf("❌ Error getting provider %s: %v\n", name, err)
			continue
		}

		fmt.Printf("\n📦 %s (v%s)\n", provider.Name(), provider.Version())

		resources := provider.Resources()
		if len(resources) == 0 {
			fmt.Println("  No resources available")
			continue
		}

		fmt.Println("  Resources:")
		for _, resource := range resources {
			fmt.Printf("    - %s: %s (%s)\n", resource.Type, resource.Name, resource.Category)
			fmt.Printf("      %s\n", resource.Description)

			if len(resource.Examples) > 0 {
				fmt.Printf("      Example: %s\n", resource.Examples[0].Name)
			}
		}
	}

	fmt.Println("\n💡 Use these resource types in your YAML configuration:")
	fmt.Println("   template: \"<provider>-<resource-type>\"")
	fmt.Println("   Example: template: \"aws-organization\"")
}

// initializeProviders registers all built-in providers with the current application version
func initializeProviders() error {
	// Register core provider with the application version
	coreProvider := core.NewCoreProviderWithVersion(version)
	if err := providers.DefaultRegistry.Register(coreProvider); err != nil {
		return fmt.Errorf("failed to register core provider '%s': %w", coreProvider.Name(), err)
	}

	// TODO: Register other built-in providers when implemented
	// awsProvider := aws.NewAWSProviderWithVersion(version)
	// if err := providers.DefaultRegistry.Register(awsProvider); err != nil {
	//     return fmt.Errorf("failed to register aws provider '%s': %w", awsProvider.Name(), err)
	// }

	return nil
}
