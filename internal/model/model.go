package model

type MapperDefinition struct {
	Name     string
	ImplName string
	Package  string
	Methods  []MappingMethod
}

type MappingMethod struct {
	Name       string
	SourceType string
	TargetType string
	Mappings   []FieldMappingRule
}

type FieldMappingRule struct {
	SourceField string
	TargetField string
	Ignore      bool
	CustomFunc  string
}
