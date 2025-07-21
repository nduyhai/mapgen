package processor

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/nduyhai/mapgen/internal/model"
)

// Processor is the interface that all processors must implement.
// Each processor is responsible for handling a specific type of directive.
type Processor interface {
	// Process processes a directive and returns a result.
	// The result can be a model.MapperDefinition, model.ValidatorDefinition, etc.
	Process(directive model.Directive) (interface{}, error)

	// Type returns the type of directive that this processor handles.
	Type() string
}

// Registry is a registry of processors.
// It maps directive types to processors.
type Registry struct {
	processors map[string]Processor
}

// NewRegistry creates a new processor registry.
func NewRegistry() *Registry {
	registry := &Registry{
		processors: make(map[string]Processor),
	}

 // Register all supported processors
	registry.Register(NewMapperProcessor())
	registry.Register(NewValidatorProcessor())
	registry.Register(NewMappingProcessor())

	return registry
}

// Register registers a processor in the registry.
func (r *Registry) Register(processor Processor) {
	r.processors[processor.Type()] = processor
}

// Get returns the processor for the given directive type.
func (r *Registry) Get(directiveType string) (Processor, bool) {
	processor, ok := r.processors[directiveType]
	return processor, ok
}

// Process processes a directive using the appropriate processor.
func (r *Registry) Process(directive model.Directive) (interface{}, error) {
	processor, ok := r.Get(directive.Type)
	if !ok {
		return nil, fmt.Errorf("no processor found for directive type: %s", directive.Type)
	}

	return processor.Process(directive)
}

// MapperProcessor is a processor for mapper directives.
type MapperProcessor struct{}

// NewMapperProcessor creates a new MapperProcessor.
func NewMapperProcessor() *MapperProcessor {
	return &MapperProcessor{}
}

// Type returns the type of directive that this processor handles.
func (p *MapperProcessor) Type() string {
	return "mapper"
}

// Process processes a mapper directive and returns a MapperDefinition.
func (p *MapperProcessor) Process(directive model.Directive) (interface{}, error) {
	// Check if the directive node is a TypeSpec
	typeSpec, ok := directive.Node.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("mapper directive must be associated with a type specification, got %T", directive.Node)
	}

	// Extract implementation name from metadata or use a default
	implName := directive.Metadata["impl"]
	if implName == "" {
		// Default to lowercase type name + "_mapper"
		implName = strings.ToLower(typeSpec.Name.Name) + "_mapper"
	}

	// Extract target file name from metadata
	targetFile := directive.Metadata["target"]

	// Get the package name from the file that contains the TypeSpec
	packageName := ""
	if file, ok := directive.Metadata["package"]; ok {
		packageName = file
	} else {
		// Default to "mapper" if package name is not provided
		packageName = "mapper"
	}

	// Create a mapper definition
	mapperDef := model.MapperDefinition{
		ImplName:   implName,
		Package:    packageName,
		TargetFile: targetFile,
		Methods:    []model.MapperMethod{},
		Imports:    []string{},
	}

	// Extract interface details if it's an interface
	if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		// Process interface methods
		for _, method := range interfaceType.Methods.List {
			if funcType, ok := method.Type.(*ast.FuncType); ok {
				// Extract method name
				methodName := ""
				if len(method.Names) > 0 {
					methodName = method.Names[0].Name
				}

				// Extract parameter types
				sourceType := ""
				if len(funcType.Params.List) > 0 {
					sourceType = exprToString(funcType.Params.List[0].Type)
				}

				// Extract return types
				targetType := ""
				if funcType.Results != nil && len(funcType.Results.List) > 0 {
					targetType = exprToString(funcType.Results.List[0].Type)
				}

				// Add method to mapper definition
				mapperDef.Methods = append(mapperDef.Methods, model.MapperMethod{
					Name:       methodName,
					SourceType: sourceType,
					TargetType: targetType,
				})
				
				// Extract package names from source and target types
				if srcPkg := extractPackage(sourceType); srcPkg != "" {
					// Check if the package is already in the imports
					found := false
					for _, imp := range mapperDef.Imports {
						if imp == srcPkg {
							found = true
							break
						}
					}
					if !found {
						mapperDef.Imports = append(mapperDef.Imports, srcPkg)
					}
				}
				
				if tgtPkg := extractPackage(targetType); tgtPkg != "" {
					// Check if the package is already in the imports
					found := false
					for _, imp := range mapperDef.Imports {
						if imp == tgtPkg {
							found = true
							break
						}
					}
					if !found {
						mapperDef.Imports = append(mapperDef.Imports, tgtPkg)
					}
				}
			}
		}
	} else {
		// For non-interface types, create standard mapper methods
		typeName := typeSpec.Name.Name
		
		// Add ToDTO method
		mapperDef.Methods = append(mapperDef.Methods, model.MapperMethod{
			Name:       "ToDTO",
			SourceType: "*" + typeName,
			TargetType: "*dto." + typeName + "DTO",
		})

		// Add ToEntity method
		mapperDef.Methods = append(mapperDef.Methods, model.MapperMethod{
			Name:       "ToEntity",
			SourceType: "*dto." + typeName + "DTO",
			TargetType: "*" + typeName,
		})
		
		// Add "dto" to imports since we're using dto package in the generated code
		mapperDef.Imports = append(mapperDef.Imports, "dto")
	}

	return mapperDef, nil
}

// ValidatorProcessor is a processor for validator directives.
type ValidatorProcessor struct{}

// NewValidatorProcessor creates a new ValidatorProcessor.
func NewValidatorProcessor() *ValidatorProcessor {
	return &ValidatorProcessor{}
}

// Type returns the type of directive that this processor handles.
func (p *ValidatorProcessor) Type() string {
	return "validator"
}

// Process processes a validator directive and returns a ValidatorDefinition.
func (p *ValidatorProcessor) Process(directive model.Directive) (interface{}, error) {
	// Check if the directive node is a TypeSpec
	typeSpec, ok := directive.Node.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("validator directive must be associated with a type specification, got %T", directive.Node)
	}

	// Extract implementation name from metadata or use a default
	implName := directive.Metadata["impl"]
	if implName == "" {
		// Default to lowercase type name + "_validator"
		implName = strings.ToLower(typeSpec.Name.Name) + "_validator"
	}

	// Create a validator definition
	validatorDef := model.ValidatorDefinition{
		ImplName: implName,
		Package:  "validator", // Default package
		Fields:   make(map[string]string),
	}

	// Extract struct fields if it's a struct
	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 {
				fieldName := field.Names[0].Name
				// For now, just add a placeholder validation rule
				validatorDef.Fields[fieldName] = "required"
			}
		}
	}

	return validatorDef, nil
}

// MappingProcessor is a processor for mapping directives.
type MappingProcessor struct{}

// NewMappingProcessor creates a new MappingProcessor.
func NewMappingProcessor() *MappingProcessor {
	return &MappingProcessor{}
}

// Type returns the type of directive that this processor handles.
func (p *MappingProcessor) Type() string {
	return "mapping"
}

// Process processes a mapping directive and returns a MappingDefinition.
func (p *MappingProcessor) Process(directive model.Directive) (interface{}, error) {
	// Create a mapping definition
	mappingDef := model.MappingDefinition{}

	// Extract mapping information from metadata
	if from, ok := directive.Metadata["from"]; ok {
		mappingDef.From = from
	}

	if to, ok := directive.Metadata["to"]; ok {
		mappingDef.To = to
	}

	if using, ok := directive.Metadata["using"]; ok {
		mappingDef.Using = using
	}

	if _, ok := directive.Metadata["ignore"]; ok {
		mappingDef.Ignore = true
	}

	return mappingDef, nil
}

// Helper function to convert an AST expression to a string representation
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// extractPackage extracts the package name from a type string
// For example, "*dto.AddressDTO" -> "dto"
func extractPackage(typeStr string) string {
	// Remove pointer prefix if present
	typeStr = strings.TrimPrefix(typeStr, "*")
	
	// Split by dot to get package and type
	parts := strings.Split(typeStr, ".")
	if len(parts) > 1 {
		return parts[0]
	}
	
	// No package qualifier
	return ""
}
