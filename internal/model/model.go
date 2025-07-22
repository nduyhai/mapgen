package model

import "go/ast"

// Directive represents a directive found in a comment.
type Directive struct {
	// Type is the type of the directive (e.g., "mapper", "validator", "mapping").
	Type string

	// Metadata is a map of metadata key-value pairs extracted from the directive.
	Metadata map[string]string

	// Node is the AST node associated with the directive.
	Node ast.Node
}

// MapperDefinition represents a mapper definition.
type MapperDefinition struct {
	Name       string
	ImplName   string
	Package    string
	Methods    []MapperMethod
	TargetFile string
	Imports    []string
}

// MapperMethod represents a method in a mapper.
type MapperMethod struct {
	Name       string
	SourceType string
	TargetType string
	Mappings   []FieldMappingRule
}

// MappingMethod represents a mapping method.
type MappingMethod struct {
	Name       string
	SourceType string
	TargetType string
	Mappings   []FieldMappingRule
}

// FieldMappingRule represents a field mapping rule.
type FieldMappingRule struct {
	From   string
	To     string
	Ignore bool
	Using  string
}

// ValidatorDefinition represents a validator definition.
type ValidatorDefinition struct {
	ImplName string
	Package  string
	Fields   map[string]string
}

// MappingDefinition represents a mapping definition.
type MappingDefinition struct {
	From       string
	To         string
	Using      string
	Ignore     bool
	MapperName string // Name of the mapper interface this mapping belongs to
	MethodName string // Name of the method this mapping belongs to
}
