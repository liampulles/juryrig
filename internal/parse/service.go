package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Service exposes functionality for parsing Go source code
type Service interface {
	ParseJuryrigged(path string) ([]MapperSpec, error)
}

// ServiceImpl implements Service
type ServiceImpl struct{}

// NewServiceImpl is a constructor
func NewServiceImpl() *ServiceImpl {
	return &ServiceImpl{}
}

// Check we implement the interface
var _ Service = &ServiceImpl{}

// ParseJuryrigged parses a "jury-rigged" filed into a spec.
// If an issues is found during parsing, an error is returned.
func (s *ServiceImpl) ParseJuryrigged(path string) ([]MapperSpec, error) {
	file, err := runParser(path)
	if err != nil {
		return nil, err
	}

	ints := findInterfaces(file)
	result := make([]MapperSpec, len(mapInts))
	for _, mapInt := range mapInts {
		mapSpec, err := extractMapSpec(&mapInt)
		if err != nil {
			return nil, err
		}
		result = append(result, *mapSpec)
	}
	return result, nil
}

type mapInt struct {
	genName string
	inter   *ast.InterfaceType
}

func filterMapInts(inters []intWithComm) []mapInt {
	var result []mapInt
	for _, inter := range inters {
		
	}
}

func extractMapSpec(mapInt *intWithComm) (*MapperSpec, error) {

}
