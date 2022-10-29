package parse

// High-level mapper types

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
	Field     string
}

type Target struct {
	Field string
}

type Function struct {
	Name       string
	Parameters map[string]string
	Results    []string
}

// Directives

type Directive interface {
	Target() Target
}

type LinkDirective struct {
	Source Source
	target Target
}

var _ Directive = &LinkDirective{}

func (ld *LinkDirective) Target() Target {
	return ld.target
}

type IgnoreDirective struct {
	target Target
}

var _ Directive = &IgnoreDirective{}

func (id *IgnoreDirective) Target() Target {
	return id.target
}

type LinkFuncDirective struct {
	Sources      []Source
	FunctionName string
	target       Target
}

var _ Directive = &LinkFuncDirective{}

func (lfd *LinkFuncDirective) Target() Target {
	return lfd.target
}
