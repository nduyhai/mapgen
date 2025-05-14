package generator

import (
	"bytes"
	"github.com/nduyhai/mapgen/internal/model"
	"os"
	"path/filepath"
	"text/template"
)

func Generate(mapper *model.MapperDefinition, outputDir string) error {
	tmpl, err := template.ParseFiles("templates/mapper_impl.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, mapper); err != nil {
		return err
	}

	filename := filepath.Join(outputDir, mapper.ImplName+".gen.go")
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	return os.WriteFile(filename, buf.Bytes(), 0644)
}
