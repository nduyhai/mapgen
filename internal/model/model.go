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
