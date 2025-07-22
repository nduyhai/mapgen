package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
				fmt.Printf("    MapperName: %s\n", mappingDef.MapperName)
				fmt.Printf("    MethodName: %s\n", mappingDef.MethodName)

				// Convert MappingDefinition to FieldMappingRule
				fieldMapping := model.FieldMappingRule{
					From:   mappingDef.From,
					To:     mappingDef.To,
					Using:  mappingDef.Using,
					Ignore: mappingDef.Ignore,
				}

				// Find the correct mapper and method to associate the mapping with
				// If MapperName is set, use it to find the mapper
				// If MethodName is set, use it to find the method
				// If neither is set, use the first method of the first mapper (fallback)

				// First, try to find the mapper by name
				mapperFound := false
				if mappingDef.MapperName != "" {
					for i, mapperDef := range mapperDefinitions {
						// Check if the mapper name matches
						if strings.EqualFold(mapperDef.ImplName, mappingDef.MapperName) {
							// If MethodName is set, find the method by name
							methodFound := false
							if mappingDef.MethodName != "" {
								for j, method := range mapperDef.Methods {
									if strings.EqualFold(method.Name, mappingDef.MethodName) {
										// Add the mapping to the method
										mapperDefinitions[i].Methods[j].Mappings = append(mapperDefinitions[i].Methods[j].Mappings, fieldMapping)
										methodFound = true
										mapperFound = true
										break
									}
								}
							}

							// If method not found but mapper found, add to the first method
							if !methodFound && len(mapperDef.Methods) > 0 {
								mapperDefinitions[i].Methods[0].Mappings = append(mapperDefinitions[i].Methods[0].Mappings, fieldMapping)
								mapperFound = true
							}

							// Break out of the mapper loop if we found the mapper
							if mapperFound {
								break
							}
						}
					}
				}

				// If mapper not found, try to find by method name
				if !mapperFound && mappingDef.MethodName != "" {
					for i, mapperDef := range mapperDefinitions {
						for j, method := range mapperDef.Methods {
							if strings.EqualFold(method.Name, mappingDef.MethodName) {
								// Add the mapping to the method
								mapperDefinitions[i].Methods[j].Mappings = append(mapperDefinitions[i].Methods[j].Mappings, fieldMapping)
								mapperFound = true
								break
							}
						}
						if mapperFound {
							break
						}
					}
				}

				// If still not found, fall back to the first method of the first mapper
				if !mapperFound && len(mapperDefinitions) > 0 && len(mapperDefinitions[0].Methods) > 0 {
					mapperDefinitions[0].Methods[0].Mappings = append(mapperDefinitions[0].Methods[0].Mappings, fieldMapping)
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
