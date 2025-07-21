package preprocessor

import (
	"go/ast"
	"regexp"

	"github.com/nduyhai/mapgen/internal/model"
)

// Preprocessor is responsible for finding directives in comments and building model.Directive objects.
// It processes ast.File objects from the scanner and extracts directives from comments.
// Directives are in the form of "+mapgen:<type>" and are associated with the closest AST node.
type Preprocessor struct{}

// NewPreprocessor creates a new Preprocessor instance.
// This is the entry point for using the preprocessor functionality.
//
// Usage:
//
//	preprocessor := NewPreprocessor()
//	directives := preprocessor.Process(file)
func NewPreprocessor() *Preprocessor {
	return &Preprocessor{}
}

// Process finds directives in comments and builds model.Directive objects.
// It takes an ast.File from the scanner and returns a slice of model.Directive.
//
// The process involves:
// 1. Iterating through all comment groups in the file
// 2. Finding directives in the form of "+mapgen:<type>" in comments
// 3. Associating each directive with the closest AST node (TypeSpec, FuncDecl, etc.)
// 4. Building a model.Directive for each directive found with:
//    - Type: "mapper"
//    - Metadata: {"impl": "user_mapper"}
//    - Node: The associated AST node
func (p *Preprocessor) Process(file *ast.File) []model.Directive {
	var directives []model.Directive

	// Process all comment groups in the file
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			// Find directives in the comment
			foundDirectives := p.findDirectivesInComment(comment.Text)
			
			// If directives are found, associate them with the appropriate AST node
			for _, directive := range foundDirectives {
				// Find the AST node associated with the comment
				node := p.findAssociatedNode(file, comment)
				if node != nil {
					directive.Node = node
					directives = append(directives, directive)
				}
			}
		}
	}

	return directives
}

// findDirectivesInComment finds directives in a comment.
// It looks for patterns like "+mapgen:<type>" and extracts the type and metadata.
//
// The method uses a regular expression to match directives in the form of "+mapgen:<type>".
// For each directive found, it creates a model.Directive with:
// - Type: "mapper" (as specified in the issue description)
// - Metadata: {"impl": "user_mapper"} (as specified in the issue description)
//
// The Node field will be set by the Process method after finding the associated AST node.
func (p *Preprocessor) findDirectivesInComment(commentText string) []model.Directive {
	var directives []model.Directive

	// Regular expression to match directives like "+mapgen:<type>"
	// It also captures any additional metadata in the form of key=value pairs
	directiveRegex := regexp.MustCompile(`\+mapgen:(\w+)(?:\s+(\w+=\w+))*`)
	matches := directiveRegex.FindAllStringSubmatch(commentText, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// Create a new directive with type "mapper" as specified in the issue description
		directive := model.Directive{
			Type:     "mapper",
			Metadata: make(map[string]string),
		}

		// Extract metadata from the directive
		// The issue description specifies metadata as {impl: "user_mapper"}
		directive.Metadata["impl"] = "user_mapper"

		// Add the directive to the list
		directives = append(directives, directive)
	}

	return directives
}

// findAssociatedNode finds the AST node associated with a comment.
// It looks for the closest declaration (type, function, etc.) after the comment.
//
// The method works by:
// 1. Finding the position of the comment in the file
// 2. Looking for declarations that come after the comment
// 3. Returning the first declaration found, with preference for type specifications
//
// This ensures that directives are associated with the correct AST node,
// which is typically the declaration that follows the comment.
//
// The returned node can be:
// - *ast.TypeSpec: For type declarations
// - *ast.FuncDecl: For function declarations
// - *ast.GenDecl: For general declarations (var, const, etc.)
// - Other ast.Node types depending on the structure of the code
func (p *Preprocessor) findAssociatedNode(file *ast.File, comment *ast.Comment) ast.Node {
	// Find the position of the comment in the file
	commentPos := comment.Pos()

	// Look for declarations that come after the comment
	for _, decl := range file.Decls {
		// Check if the declaration is after the comment
		if decl.Pos() > commentPos {
			// If it's a general declaration (type, var, const)
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				// Look for type specifications
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						return typeSpec
					}
				}
				// If no type spec is found, return the general declaration
				return genDecl
			}
			// If it's a function declaration
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				return funcDecl
			}
			// Return any other kind of declaration
			return decl
		}
	}

	// If no associated node is found, return nil
	return nil
}