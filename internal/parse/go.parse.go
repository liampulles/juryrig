package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type intWithComm struct {
	commGrp *ast.CommentGroup
	inter   *ast.InterfaceType
}

func runParser(path string) (*ast.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("parser failed: %w", err)
	}
	return f, nil
}

func findInterfaces(file *ast.File) []intWithComm {
	var result []intWithComm
	for _, decl := range file.Decls {
		intTypes := declToInterfaceTypes(decl)
		result = append(result, intTypes...)
	}
	return result
}

func declToInterfaceTypes(decl ast.Decl) []intWithComm {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil
	}

	if genDecl.Tok != token.TYPE {
		return nil
	}

	var result []intWithComm
	for _, spec := range genDecl.Specs {
		intType := specToInterfaceType(spec)
		if intType == nil {
			continue
		}
		result = append(result, *intType)
	}
	return result
}

func specToInterfaceType(spec ast.Spec) *intWithComm {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	intType, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return nil
	}

	return &intWithComm{
		commGrp: typeSpec.Comment,
		inter:   intType,
	}
}
