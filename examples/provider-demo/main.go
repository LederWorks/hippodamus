package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/LederWorks/hippodamus/pkg/providers"
	"github.com/LederWorks/hippodamus/providers/aws"
)

func main() {
	// Create a provider registry
	registry := providers.NewRegistry()

	// Register the AWS provider
	awsProvider := aws.NewAWSProvider()
	if err := registry.Register(awsProvider); err != nil {
		log.Fatalf("Failed to register AWS provider: %v", err)
	}

	fmt.Println("🚀 Hippodamus Provider System Demo")
	fmt.Println("===================================")

	// List all registered providers
	fmt.Printf("📋 Registered providers: %v\n", registry.List())

	// Get AWS provider and show its resources
	provider, err := registry.Get("aws")
	if err != nil {
		log.Fatalf("Failed to get AWS provider: %v", err)
	}

	fmt.Printf("🔧 AWS Provider v%s\n", provider.Version())
	fmt.Println("📦 Available resources:")

	resources := provider.Resources()
	for _, resource := range resources {
		fmt.Printf("  - %s: %s (%s)\n", resource.Type, resource.Name, resource.Category)
	}

	// Test AWS Organization generation
	fmt.Println("\n🏗️  Generating AWS Organization template...")
	orgParams := map[string]interface{}{
		"orgName":             "Demo Organization",
		"managementAccountId": "123456789012",
		"fillColor":           "#FFF8E1",
		"strokeColor":         "#FF9900",
	}

	element, err := provider.GenerateTemplate("organization", orgParams)
	if err != nil {
		log.Fatalf("Failed to generate organization template: %v", err)
	}

	fmt.Printf("✅ Generated element: %s (type: %s)\n", element.Name, element.Type)
	fmt.Printf("   Label: %s\n", element.Properties.Label)
	fmt.Printf("   Dimensions: %.0fx%.0f\n", element.Properties.Width, element.Properties.Height)
	fmt.Printf("   Colors: %s / %s\n", element.Style.FillColor, element.Style.StrokeColor)

	// Test VPC generation
	fmt.Println("\n🌐 Generating AWS VPC template...")
	vpcParams := map[string]interface{}{
		"vpcName":   "production-vpc",
		"cidrBlock": "10.0.0.0/16",
		"region":    "us-west-2",
	}

	vpcElement, err := provider.GenerateTemplate("vpc", vpcParams)
	if err != nil {
		log.Fatalf("Failed to generate VPC template: %v", err)
	}

	fmt.Printf("✅ Generated element: %s (type: %s)\n", vpcElement.Name, vpcElement.Type)
	fmt.Printf("   Label: %s\n", vpcElement.Properties.Label)

	// Show resource schema
	fmt.Println("\n📋 AWS Organization Schema:")
	schema, err := provider.GetSchema("organization")
	if err != nil {
		log.Fatalf("Failed to get schema: %v", err)
	}

	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(schemaJSON))

	// Test validation
	fmt.Println("\n🔍 Testing validation...")

	// Valid case
	err = provider.Validate("organization", orgParams)
	if err != nil {
		fmt.Printf("❌ Validation failed (unexpected): %v\n", err)
	} else {
		fmt.Printf("✅ Valid organization parameters\n")
	}

	// Invalid case
	invalidParams := map[string]interface{}{
		"orgName": "Test Org",
		// Missing managementAccountId
	}
	err = provider.Validate("organization", invalidParams)
	if err != nil {
		fmt.Printf("✅ Validation correctly failed: %v\n", err)
	} else {
		fmt.Printf("❌ Validation should have failed\n")
	}

	fmt.Println("\n🎉 Provider system demo completed successfully!")
}
