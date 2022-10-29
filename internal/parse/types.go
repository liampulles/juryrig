package parse

// High-level mapper types

type JuryrigSpec struct {
	Package string
	Mappers []Mapper
}

type Mapper struct {
	Name            string
	MapperFunctions []MapperFunction
}

type MapperFunction struct {
	Function   Function
	Directives []Directive
}

// Common types

type Source struct {
	Parameter string
	// Optional
	Field string
}

type Target struct {
	Field string
}

type Function struct {
	Name       string
	Parameters map[string]string
	Results    []string
}

// Directive indicates how to handle a given target field on a struct.
type Directive interface{}

type LinkDirective struct {
	Source Source
	Target Target
}

var _ Directive = &LinkDirective{} //nolint:exhaustruct

type IgnoreDirective struct {
	Target Target
}

var _ Directive = &IgnoreDirective{} //nolint:exhaustruct

type LinkFuncDirective struct {
	Sources      []Source
	FunctionName string
	Target       Target
}

var _ Directive = &LinkFuncDirective{} //nolint:exhaustruct
