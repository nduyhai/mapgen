package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nduyhai/mapgen/internal/generator"
	"github.com/nduyhai/mapgen/internal/model"
	"github.com/nduyhai/mapgen/internal/preprocessor"
	"github.com/nduyhai/mapgen/internal/processor"
	"github.com/nduyhai/mapgen/internal/scanner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mapgen <file_or_directory_path>")
		os.Exit(1)
	}
	// Create a new scanner
	s := scanner.NewScanner()

	// Get the path from command line arguments
	path := os.Args[1]

	// Check if the path is a file or directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Error accessing path %s: %v", path, err)
	}

	if fileInfo.IsDir() {
		// Parse and type check directory
		fmt.Printf("Parsing and type checking directory: %s\n", path)
		pkgs, err := s.ParsePackage(path)
		if err != nil {
			log.Fatalf("Error processing directory: %v", err)
		}

		// Print the packages found
		fmt.Printf("Found %d packages:\n", len(pkgs))
		for _, pkg := range pkgs {
			fmt.Printf("- %s (%s)\n", pkg.Name(), pkg.Path())
		}

	} else {
		// Parse and type check file
		fmt.Printf("Parsing and type checking file: %s\n", path)
		file, typeInfo, err := s.TypeCheckFile(path)
		if err != nil {
			log.Fatalf("Error processing file: %v", err)
		}

		// Print some basic information about the file
		fmt.Printf("Package: %s\n", file.Name.Name)
		fmt.Printf("Number of imports: %d\n", len(file.Imports))
		fmt.Printf("Number of declarations: %d\n", len(file.Decls))
		fmt.Printf("Number of comments: %d\n", len(file.Comments))

		// Print type information
		fmt.Printf("\nType Information:\n")
		fmt.Printf("- Number of types: %d\n", len(typeInfo.Types))
		fmt.Printf("- Number of definitions: %d\n", len(typeInfo.Defs))
		fmt.Printf("- Number of uses: %d\n", len(typeInfo.Uses))

		// Create a preprocessor and process the file
		p := preprocessor.NewPreprocessor()
		directives := p.Process(file)

		// Print information about the directives found
		fmt.Printf("\nDirectives:\n")
		fmt.Printf("- Number of directives: %d\n", len(directives))
		for i, directive := range directives {
			fmt.Printf("  Directive %d:\n", i+1)
			fmt.Printf("  - Type: %s\n", directive.Type)
			fmt.Printf("  - Metadata: %v\n", directive.Metadata)
			fmt.Printf("  - Node Type: %T\n", directive.Node)
		}

		// Create a processor registry
		registry := processor.NewRegistry()

		// Group directives by type
		mapperDirectives := []model.Directive{}
		mappingDirectives := []model.Directive{}
		validatorDirectives := []model.Directive{}

		for _, directive := range directives {
			switch directive.Type {
			case "mapper":
				mapperDirectives = append(mapperDirectives, directive)
			case "mapping":
				mappingDirectives = append(mappingDirectives, directive)
			case "validator":
				validatorDirectives = append(validatorDirectives, directive)
			}
		}

		fmt.Printf("\nProcessed Results:\n")

		// Process mapper directives first
		fmt.Printf("\nMapper Directives:\n")
		mapperDefinitions := []model.MapperDefinition{}

		for i, directive := range mapperDirectives {
			result, err := registry.Process(directive)
			if err != nil {
				fmt.Printf("  Error processing mapper directive %d: %v\n", i+1, err)
				continue
			}

			if mapperDef, ok := result.(model.MapperDefinition); ok {
				fmt.Printf("  Mapper %d:\n", i+1)
				fmt.Printf("    ImplName: %s\n", mapperDef.ImplName)
				fmt.Printf("    Package: %s\n", mapperDef.Package)
				fmt.Printf("    TargetFile: %s\n", mapperDef.TargetFile)
				fmt.Printf("    Methods:\n")
				for j, method := range mapperDef.Methods {
					fmt.Printf("      Method %d:\n", j+1)
					fmt.Printf("        Name: %s\n", method.Name)
					fmt.Printf("        SourceType: %s\n", method.SourceType)
					fmt.Printf("        TargetType: %s\n", method.TargetType)
				}

				// Store the mapper definition
				mapperDefinitions = append(mapperDefinitions, mapperDef)
			}
		}

		// Process mapping directives
		fmt.Printf("\nMapping Directives:\n")
		for i, directive := range mappingDirectives {
			result, err := registry.Process(directive)
			if err != nil {
				fmt.Printf("  Error processing mapping directive %d: %v\n", i+1, err)
				continue
			}

			if mappingDef, ok := result.(model.MappingDefinition); ok {
				fmt.Printf("  Mapping %d:\n", i+1)
				fmt.Printf("    From: %s\n", mappingDef.From)
				fmt.Printf("    To: %s\n", mappingDef.To)
				fmt.Printf("    Using: %s\n", mappingDef.Using)
				fmt.Printf("    Ignore: %v\n", mappingDef.Ignore)

				// For now, we'll associate the mapping with the first method of the first mapper definition
				// In a real implementation, we would need to find the correct method based on the AST
				if len(mapperDefinitions) > 0 && len(mapperDefinitions[0].Methods) > 0 {
					// Add the mapping to the first method of the first mapper definition
					mapperDefinitions[0].Methods[0].Mappings = append(mapperDefinitions[0].Methods[0].Mappings, mappingDef)
				}
			}
		}

		// Process validator directives
		fmt.Printf("\nValidator Directives:\n")
		for i, directive := range validatorDirectives {
			result, err := registry.Process(directive)
			if err != nil {
				fmt.Printf("  Error processing validator directive %d: %v\n", i+1, err)
				continue
			}

			if validatorDef, ok := result.(model.ValidatorDefinition); ok {
				fmt.Printf("  Validator %d:\n", i+1)
				fmt.Printf("    ImplName: %s\n", validatorDef.ImplName)
				fmt.Printf("    Package: %s\n", validatorDef.Package)
				fmt.Printf("    Fields: %v\n", validatorDef.Fields)
			}
		}

		// Generate code for mapper definitions
		for _, mapperDef := range mapperDefinitions {
			g := generator.NewGenerator(path)
			err := g.Generate(mapperDef)
			if err != nil {
				fmt.Printf("  Error generating code: %v\n", err)
			}
		}
	}
}
