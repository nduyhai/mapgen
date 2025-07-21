package main

import (
	"fmt"
	"log"
	"os"

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
	}
}
