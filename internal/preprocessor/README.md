# Preprocessor Package

The preprocessor package is responsible for finding directives in comments and building model.Directive objects. It processes ast.File objects from the scanner and extracts directives from comments.

## Overview

The preprocessor works by:
1. Iterating through all comment groups in the file
2. Finding directives in the form of "+mapgen:<type>" in comments
3. Associating each directive with the closest AST node (TypeSpec, FuncDecl, etc.)
4. Building a model.Directive for each directive found with:
   - Type: "mapper"
   - Metadata: {"impl": "user_mapper"}
   - Node: The associated AST node

## Usage

```go
// Create a scanner and parse a file
scanner := scanner.NewScanner()
file, _, err := scanner.TypeCheckFile(filePath)
if err != nil {
    // Handle error
}

// Create a preprocessor and process the file
preprocessor := preprocessor.NewPreprocessor()
directives := preprocessor.Process(file)

// Use the directives
for _, directive := range directives {
    fmt.Printf("Type: %s\n", directive.Type)
    fmt.Printf("Metadata: %v\n", directive.Metadata)
    fmt.Printf("Node Type: %T\n", directive.Node)
}
```

## Directive Format

Directives are comments in the form of "+mapgen:<type>", where <type> is the type of the directive. For example:

```go
// +mapgen:mapper
type User struct {
    ID   int
    Name string
}
```

The preprocessor will find this directive and associate it with the User struct.

## Implementation Details

The preprocessor uses the following methods:

- `NewPreprocessor()`: Creates a new Preprocessor instance
- `Process(file *ast.File)`: Processes the file and returns a slice of model.Directive
- `findDirectivesInComment(commentText string)`: Finds directives in a comment
- `findAssociatedNode(file *ast.File, comment *ast.Comment)`: Finds the AST node associated with a comment