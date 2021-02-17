package parse

import (
	"go/ast"
)

type MapperSpec struct {
	Name  string
	Funcs []MapperFuncSpec
}

type MapperFuncSpec struct {
	Name            string
	AnnotationSpecs []interface{}
	FuncType        ast.FuncType
}

type LinkSpec struct {
	Source string
	Target string
}

type IgnoreSpec struct {
	Target string
}

type LinkFuncSpec struct {
	Sources  []string
	FuncName string
	Target   string
}
