package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const (
	juryRigTag       = "// +juryrig:"
	juryRigMapperTag = juryRigTag + "mapper:"
)

// ParseFileWithMapperDefs parses Go source into Mapper definitions
func ParseFileWithMapperDefs(filename string) ([]Mapper, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse file %s: %w", filename, err)
	}

	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}
	result, err := parseMappers(fset, body, astFile)
	if err != nil {
		return nil, fmt.Errorf("could not parse ast: %w", err)
	}
	return result, nil
}

func parseMappers(fset *token.FileSet, body []byte, astFile *ast.File) ([]Mapper, error) {
	var result []Mapper
	for _, decl := range astFile.Decls {
		if !isJuryRigMapperDecl(decl) {
			continue
		}
		raw, err := extractRawMapperInfo(fset, astFile, body, decl)
		if err != nil {
			return nil, fmt.Errorf("could not extract mapper: %w", err)
		}
		fmt.Println(raw)
	}
	return result, nil
}

type rawMapperInfo struct {
	name          string
	topJrComments []string
	fns           []*rawMapperFuncInfo
}

type rawMapperFuncInfo struct {
	name       string
	parameters map[string]string
	results    []string
	jrComments []string
}

func isJuryRigMapperDecl(decl ast.Decl) bool {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return false
	}
	return isJuryRigMapperCommentGroup(genDecl.Doc)
}

func extractRawMapperInfo(fset *token.FileSet, astFile *ast.File, body []byte, decl ast.Decl) (*rawMapperInfo, error) {
	result := &rawMapperInfo{}
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil, fmt.Errorf("mapper tagged type is not a declaration %s",
			locationInfo(fset, decl))
	}

	result.topJrComments = filterJuryRigComments(genDecl.Doc.List, juryRigTag)

	if len(genDecl.Specs) != 1 {
		return nil, fmt.Errorf("expect mapper declaration to have 1 spec, but has %d",
			len(genDecl.Specs))
	}
	spec := genDecl.Specs[0]
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("mapper spec is not a type %s",
			locationInfo(fset, genDecl))
	}

	result.name = typeSpec.Name.Name

	intSpec, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return nil, fmt.Errorf("mapper type is not an interface %s",
			locationInfo(fset, typeSpec))
	}

	fnInfos, err := extractRawMapperFuncInfos(fset, astFile, body, intSpec.Methods.List)
	if err != nil {
		return nil, fmt.Errorf("could not extract methods for mapper: %w", err)
	}
	result.fns = fnInfos

	return result, nil
}

func extractRawMapperFuncInfos(fset *token.FileSet, astFile *ast.File, body []byte, fields []*ast.Field) ([]*rawMapperFuncInfo, error) {
	var result []*rawMapperFuncInfo
	for i, field := range fields {
		fnInfo, err := extractRawMapperFuncInfo(fset, astFile, body, field)
		if err != nil {
			return nil, fmt.Errorf("could not extract field %d: %w", i, err)
		}
		result = append(result, fnInfo)
	}
	return result, nil
}

func extractRawMapperFuncInfo(fset *token.FileSet, astFile *ast.File, body []byte, field *ast.Field) (*rawMapperFuncInfo, error) {
	result := &rawMapperFuncInfo{}
	if len(field.Names) != 1 {
		return nil, fmt.Errorf("expected field to have 1 name, but has %d %s",
			len(field.Names), locationInfo(fset, field))
	}
	result.name = field.Names[0].Name

	funcType, ok := field.Type.(*ast.FuncType)
	if !ok {
		return nil, fmt.Errorf("field type is not a function %s",
			locationInfo(fset, field))
	}

	params, err := extractFuncParamaters(fset, astFile, body, funcType)
	if err != nil {
		return nil, fmt.Errorf("could not extract func parameters: %w", err)
	}
	result.parameters = params

	result.results = extractFuncResults(fset, astFile, body, funcType)
	result.jrComments = filterJuryRigComments(field.Doc.List, juryRigTag)

	return result, nil
}

func extractFuncParamaters(fset *token.FileSet, astFile *ast.File, body []byte, fn *ast.FuncType) (map[string]string, error) {
	result := make(map[string]string)
	for i, paramField := range fn.Params.List {
		if len(paramField.Names) != 1 {
			return nil, fmt.Errorf("paramater %d does not have one name %s",
				i, locationInfo(fset, paramField))
		}
		name := paramField.Names[0].Name
		typ := readAsString(fset, astFile, body, paramField.Type)
		result[name] = typ
	}
	return result, nil
}

func extractFuncResults(fset *token.FileSet, astFile *ast.File, body []byte, fn *ast.FuncType) []string {
	var result []string
	for _, resultField := range fn.Results.List {
		typ := readAsString(fset, astFile, body, resultField.Type)
		result = append(result, typ)
	}
	return result
}

func extractMapperComment(jrComments []string) (string, error) {
	if len(jrComments) != 1 {
		return "", fmt.Errorf("expecting mapper declaration to have exactly 1 juryrig tag, but has %d",
			len(jrComments))
	}
	jrComment := jrComments[0]
	if !strings.HasPrefix(jrComment, juryRigMapperTag) {
		return "", fmt.Errorf("expecting mapper declaration to have juryrig tag of the form %s..., but was %s",
			juryRigMapperTag, jrComment)
	}
	return jrComment, nil
}

func isJuryRigMapperCommentGroup(commentGroup *ast.CommentGroup) bool {
	for _, cmt := range commentGroup.List {
		if isJuryRigComment(cmt, juryRigMapperTag) {
			return true
		}
	}
	return false
}

func filterJuryRigComments(comments []*ast.Comment, tag string) []string {
	var result []string
	for _, cmt := range comments {
		if isJuryRigComment(cmt, tag) {
			result = append(result, strings.TrimSpace(cmt.Text))
		}
	}
	return result
}

func isJuryRigComment(comment *ast.Comment, tag string) bool {
	return strings.HasPrefix(strings.TrimSpace(comment.Text), tag)
}

func locationInfo(fset *token.FileSet, node ast.Node) string {
	position := fset.Position(node.Pos())
	return fmt.Sprintf("[file %s, line %d]", position.Filename, position.Line)
}

func readAsString(fset *token.FileSet, astFile *ast.File, body []byte, node ast.Node) string {
	offset := astFile.Pos()
	return string(body[node.Pos()-offset : node.End()-offset])
}
