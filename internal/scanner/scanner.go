package scanner

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

// ModuleImporter is a custom importer that can resolve packages in the local module
type ModuleImporter struct {
	moduleRoot  string
	stdImporter types.Importer
}

// Import implements the types.Importer interface
func (mi *ModuleImporter) Import(path string) (*types.Package, error) {
	// First try the standard importer for standard library packages
	pkg, err := mi.stdImporter.Import(path)
	if err == nil {
		return pkg, nil
	}

	// Check if the import path starts with the module name from go.mod
	if strings.HasPrefix(path, "github.com/nduyhai/mapgen/") {
		// Extract the relative path within the module
		relPath := strings.TrimPrefix(path, "github.com/nduyhai/mapgen/")

		// Construct the local filesystem path
		localPath := filepath.Join(mi.moduleRoot, relPath)

		// Parse the package
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, localPath, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse package %s: %w", path, err)
		}

		if len(pkgs) == 0 {
			return nil, fmt.Errorf("no packages found at %s", localPath)
		}

		// Use the first package found
		var astPkg *ast.Package
		for _, pkg := range pkgs {
			astPkg = pkg
			break
		}

		// Create a new types.Package
		typesPkg := types.NewPackage(path, astPkg.Name)

		// Convert package.Files map to a slice of *ast.File
		var files []*ast.File
		for _, file := range astPkg.Files {
			files = append(files, file)
		}

		// Type check the package
		typeConfig := &types.Config{
			Importer: mi,                 // Use this importer for recursive imports
			Error:    func(err error) {}, // Silently collect errors
		}

		typeInfo := &types.Info{}
		err = types.NewChecker(typeConfig, fset, typesPkg, typeInfo).Files(files)
		if err != nil {
			return nil, fmt.Errorf("type checking error for package %s: %w", path, err)
		}

		return typesPkg, nil
	}

	// Fall back to the standard importer for other packages
	return mi.stdImporter.Import(path)
}

// Scanner is responsible for parsing Go source files into ASTs and performing type checking.
type Scanner struct {
	// FileSet provides position information for AST nodes
	fset *token.FileSet
	// Module root directory
	moduleRoot string
}

// NewScanner creates a new Scanner instance.
// It initializes the token.FileSet for position information and determines the module root.
func NewScanner() *Scanner {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Warning: could not determine current directory: %v\n", err)
		cwd = "."
	}

	// Find the module root (directory containing go.mod)
	moduleRoot := findModuleRoot(cwd)

	return &Scanner{
		fset:       token.NewFileSet(),
		moduleRoot: moduleRoot,
	}
}

// findModuleRoot finds the directory containing the go.mod file
// by traversing up the directory tree from the given path.
func findModuleRoot(dir string) string {
	for {
		// Check if go.mod exists in the current directory
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root directory without finding go.mod
			fmt.Println("Warning: could not find go.mod file, module-relative imports may not work")
			return ""
		}
		dir = parent
	}
}

// ParseFile parses a single Go source file into an AST
func (s *Scanner) ParseFile(filePath string) (*ast.File, error) {
	// Check if a file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Parse the file
	file, err := parser.ParseFile(s.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	return file, nil
}

// ParseDir parses all Go source files in a directory and converts them to types.Package.
func (s *Scanner) ParseDir(dirPath string) ([]*types.Package, error) {
	// Check if a directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dirPath)
	}

	// Parse the directory to get AST packages
	astPkgs, err := parser.ParseDir(s.fset, dirPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse directory %s: %w", dirPath, err)
	}

	// Convert ast.Package to types.Package
	var typesPkgs []*types.Package
	for name, astPkg := range astPkgs {
		// Create a new types.Package
		typesPkg := types.NewPackage(dirPath, name)

		// Convert package.Files map to a slice of *ast.File
		var files []*ast.File
		for _, file := range astPkg.Files {
			files = append(files, file)
		}

		// Create type info for this check
		typeInfo := &types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
			Scopes:     make(map[ast.Node]*types.Scope),
		}

		// Create a custom module-aware importer
		customImporter := &ModuleImporter{
			moduleRoot:  s.moduleRoot,
			stdImporter: importer.Default(),
		}

		// Create type config with the custom importer
		typeConfig := &types.Config{
			Importer: customImporter,
			Error:    func(err error) {}, // Silently collect errors
		}

		err = types.NewChecker(typeConfig, s.fset, typesPkg, typeInfo).Files(files)
		if err != nil {
			// Continue with other packages even if one fails
			fmt.Printf("Warning: type checking error in package %s: %v\n", name, err)
		}

		typesPkgs = append(typesPkgs, typesPkg)
	}

	return typesPkgs, nil
}

// ParsePackage parses all Go source files in a package and converts them to types.Package.
// This method replaces the deprecated ast.Package with types.Package as recommended.
func (s *Scanner) ParsePackage(pkgPath string) ([]*types.Package, error) {
	// Resolve the absolute path
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path for %s: %w", pkgPath, err)
	}

	return s.ParseDir(absPath)
}

// GetFileSet returns the token.FileSet used by the scanner
func (s *Scanner) GetFileSet() *token.FileSet {
	return s.fset
}

// TypeCheckFile parses and type checks a single Go source file.
// It returns the AST, type information, and any error that occurred.
func (s *Scanner) TypeCheckFile(filePath string) (*ast.File, *types.Info, error) {
	// Parse the file first
	file, err := s.ParseFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	// Create a new package to type check the file
	pkg := types.NewPackage(filepath.Dir(filePath), file.Name.Name)

	// Create type info for this check
	typeInfo := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	// Create a custom module-aware importer
	customImporter := &ModuleImporter{
		moduleRoot:  s.moduleRoot,
		stdImporter: importer.Default(),
	}

	// Create type config with the custom importer
	typeConfig := &types.Config{
		Importer: customImporter,
		Error:    func(err error) {}, // Silently collect errors
	}

	// Type checks the file
	err = types.NewChecker(typeConfig, s.fset, pkg, typeInfo).Files([]*ast.File{file})
	if err != nil {
		return file, typeInfo, fmt.Errorf("type checking error: %w", err)
	}

	return file, typeInfo, nil
}
