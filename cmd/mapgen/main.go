package main

import (
	"flag"
	"github.com/nduyhai/mapgen/internal/generator"
	"github.com/nduyhai/mapgen/internal/parser"
	"log"
)

func main() {
	input := flag.String("input", "./example/mapper", "Directory to search for mapper interfaces")
	output := flag.String("output", "./example/mapper_gen", "Directory to output generated code")
	flag.Parse()

	mappers, err := parser.ParseDir(*input)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, mapper := range mappers {
		if err := generator.Generate(mapper, *output); err != nil {
			log.Fatal(err)
			return
		}
	}
}
