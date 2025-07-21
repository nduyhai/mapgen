package model

import (
	"go/ast"
)

// Directive represents a code generation directive found in comments.
// It contains information about the type of directive, its metadata,
// and the AST node it's associated with.
type Directive struct {
	// Type is the type of the directive (e.g., "mapper")
	Type string

	// Metadata contains additional information about the directive
	// (e.g., {"impl": "user_mapper"})
	Metadata map[string]string

	// Node is the AST node associated with the directive
	// (e.g., *ast.TypeSpec or similar)
	Node ast.Node
}

// MappingDefinition represents a field mapping in a mapper method
type MappingDefinition struct {
	// From is the source field name
	From string

	// To is the target field name
	To string

	// Using is the function to use for the mapping
	Using string

	// Ignore indicates whether the field should be ignored in the mapping
	Ignore bool
}

// MapperMethod represents a method in a mapper interface
type MapperMethod struct {
	// Name is the name of the method
	Name string

	// SourceType is the type of the source parameter
	SourceType string

	// TargetType is the type of the target parameter or return value
	TargetType string

	// Mappings is a slice of field mappings for this method
	Mappings []MappingDefinition
}

// MapperDefinition represents a definition of a mapper interface
type MapperDefinition struct {
	// ImplName is the name of the implementation (e.g., "user_mapper")
	ImplName string

	// Package is the package of the implementation (e.g., "mapper")
	Package string

	// TargetFile is the name of the target file for the generated code
	// If not specified, the source file name with "_mapper" appended will be used
	TargetFile string

	// Methods is a slice of methods in the mapper interface
	Methods []MapperMethod

	// Imports is a slice of packages that need to be imported in the generated code
	Imports []string
}

// ValidatorDefinition represents a definition of a validator interface
// This is a placeholder for future implementation
type ValidatorDefinition struct {
	// ImplName is the name of the implementation
	ImplName string

	// Package is the package of the implementation
	Package string

	// Fields is a map of field names to validation rules
	Fields map[string]string
}
