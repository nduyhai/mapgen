# Model Package

The model package contains data structures used throughout the mapgen project. It defines the core models that represent directives found in comments and their associated AST nodes.

## Directive

The `Directive` struct represents a code generation directive found in comments. It contains information about the type of directive, its metadata, and the AST node it's associated with.

### Structure

```go
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
```

### Fields

- **Type**: A string representing the type of the directive. In the current implementation, this is always "mapper".
- **Metadata**: A map of string keys to string values containing additional information about the directive. In the current implementation, this is always `{"impl": "user_mapper"}`.
- **Node**: The AST node associated with the directive. This can be:
  - `*ast.TypeSpec`: For type declarations
  - `*ast.FuncDecl`: For function declarations
  - `*ast.GenDecl`: For general declarations (var, const, etc.)
  - Other `ast.Node` types depending on the structure of the code

### Usage

The `Directive` struct is primarily used by the preprocessor to represent directives found in comments. The preprocessor builds a slice of `Directive` objects and returns them for further processing.

```go
// Create a preprocessor and process a file
preprocessor := preprocessor.NewPreprocessor()
directives := preprocessor.Process(file)

// Use the directives
for _, directive := range directives {
    // Access the directive's fields
    fmt.Printf("Type: %s\n", directive.Type)
    fmt.Printf("Metadata: %v\n", directive.Metadata)
    fmt.Printf("Node Type: %T\n", directive.Node)
}
```