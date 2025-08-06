# Template Dependency System Proposal

## Problem Statement

Currently, Azure DevOps resources have implicit hierarchical dependencies:
- **Projects** depend on **Organizations** 
- **Repositories** depend on **Projects**
- **Pipelines** depend on **Repositories**
- **Service Connections** can depend on **Projects** or **Repositories**

However, our templates don't model these dependencies explicitly, leading to:
- No validation of required parent resources
- Difficulty building dependency trees
- Complex manual management of hierarchical relationships
- No automated parent resolution

## Proposed Solution

### 1. Template Dependency Declaration

Add explicit `dependencies` section to templates:

```yaml
# azuredevops-project.yaml
name: "azuredevops-project"
description: "Azure DevOps Project container"
version: "1.0"

# NEW: Explicit dependency declaration
dependencies:
  - name: "organization"
    type: "azuredevops-organization"
    required: true
    description: "Parent Azure DevOps organization"
    relationship: "parent"

parameters:
  - name: "projectName"
    type: "string"
    default: "Project"
    description: "Project name"
  # ... other parameters
```

### 2. Enhanced Schema Types

Extend the `Template` struct:

```go
type Template struct {
    Name         string       `yaml:"name" json:"name"`
    Description  string       `yaml:"description,omitempty" json:"description,omitempty"`
    Version      string       `yaml:"version,omitempty" json:"version,omitempty"`
    Dependencies []Dependency `yaml:"dependencies,omitempty" json:"dependencies,omitempty"` // NEW
    Parameters   []Parameter  `yaml:"parameters,omitempty" json:"parameters,omitempty"`
    Elements     []Element    `yaml:"elements" json:"elements"`
}

type Dependency struct {
    Name         string `yaml:"name" json:"name"`                           // Logical name for the dependency
    Type         string `yaml:"type" json:"type"`                           // Template type required
    Required     bool   `yaml:"required,omitempty" json:"required,omitempty"` // Whether dependency is mandatory
    Description  string `yaml:"description,omitempty" json:"description,omitempty"`
    Relationship string `yaml:"relationship" json:"relationship"`           // "parent", "peer", "child"
    Multiple     bool   `yaml:"multiple,omitempty" json:"multiple,omitempty"` // Allow multiple instances
}
```

### 3. Dependency Hierarchy

**Azure DevOps Resource Hierarchy:**
```
Organization (root)
└── Project (depends on: organization)
    ├── Repository (depends on: project)
    │   ├── Pipeline (depends on: repository)
    │   └── Service Connection (depends on: repository OR project)
    ├── Library (depends on: project)
    └── Environment (depends on: project)
```

### 4. Template Updates

Each template would declare its dependencies:

#### azuredevops-project.yaml
```yaml
dependencies:
  - name: "organization"
    type: "azuredevops-organization"
    required: true
    relationship: "parent"
```

#### azuredevops-repository.yaml
```yaml
dependencies:
  - name: "project"
    type: "azuredevops-project"
    required: true
    relationship: "parent"
```

#### azuredevops-pipeline.yaml
```yaml
dependencies:
  - name: "repository"
    type: "azuredevops-repository"
    required: true
    relationship: "parent"
  - name: "project"
    type: "azuredevops-project"
    required: false
    relationship: "ancestor"
    description: "Resolved automatically through repository->project chain"
```

### 5. Configuration Validation

With dependencies declared, we can validate configurations:

```yaml
# INVALID - Missing organization parent
diagram:
  pages:
    - elements:
        - template: "azuredevops-project"  # ERROR: No organization parent found
```

```yaml
# VALID - Proper hierarchy
diagram:
  pages:
    - elements:
        - template: "azuredevops-organization"
          children:
            - template: "azuredevops-project"  # OK: Has organization parent
```

### 6. Enhanced Processing

The template processor would:
1. **Validate dependencies** before processing
2. **Build dependency tree** for proper ordering
3. **Auto-resolve ancestor relationships** (project->organization through repository)
4. **Provide helpful error messages** for missing dependencies

## Benefits

### 1. **Explicit Modeling**
- Clear declaration of resource relationships
- Self-documenting templates
- Better understanding of Azure DevOps structure

### 2. **Validation & Safety**
- Catch configuration errors early
- Prevent invalid resource combinations
- Ensure proper hierarchy

### 3. **Simplified Configuration**
- Auto-resolve parent relationships
- Reduce manual hierarchy management
- Better error messages

### 4. **Enhanced Tooling**
- Dependency tree visualization
- Automated parent resolution
- Template compatibility checking

## Implementation Steps

1. **Extend schema types** - Add `Dependency` struct and update `Template`
2. **Update template files** - Add dependency declarations
3. **Enhance template processor** - Add dependency validation
4. **Update configuration** - Test with new dependency system
5. **Add tooling** - Dependency tree visualization, validation reports

## Example Enhanced Templates

This would make the Azure DevOps resource model much more robust and easier to understand while providing better validation and tooling support.
