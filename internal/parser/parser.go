package parser

import (
	"github.com/nduyhai/mapgen/internal/model"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func ParseDir(dir string) ([]*model.MapperDefinition, error) {
	var mappers []*model.MapperDefinition

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return nil
		}
		fileSet := token.NewFileSet()
		node, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, decl := range node.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				_, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				if gen.Doc == nil {
					continue
				}
				for _, comment := range gen.Doc.List {
					if strings.HasPrefix(comment.Text, "// +mapgen:mapper") {
						mapper := &model.MapperDefinition{
							Name:     typeSpec.Name.Name,
							ImplName: typeSpec.Name.Name,
							Package:  node.Name.Name,
						}
						// TODO: parse method mappings
						mappers = append(mappers, mapper)
					}
				}
			}
		}
		return nil
	})

	return mappers, err
}
