package parse

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type rawMapperInfo struct {
	name          string
	topJrComments []string
	fns           []rawMapperFuncInfo
}

type rawMapperFuncInfo struct {
	name       string
	parameters []Parameter
	result     string
	jrComments []string
}

const (
	juryRigTag = "// +juryrig:"
)

var (
	ErrUnexpectedAST = errors.New("unexpected AST")
	ErrSpec          = errors.New("specification error")
)

// Just extract the most basic raw details from the files. Keep the
// ast stuff here basically.
func extractRaw(filename string) (string, []rawMapperInfo, error) {
	// Read ast and files
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)

	if err != nil {
		return "", nil, fmt.Errorf("could not parse file %s: %w", filename, err)
	}

	body, err := os.ReadFile(filename)

	if err != nil {
		return "", nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}

	// Parse
	pkg := astFile.Name.Name

	result, err := parseMappers(fset, body, astFile)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse ast: %w", err)
	}

	return pkg, result, nil
}

func parseMappers(fset *token.FileSet, body []byte, astFile *ast.File) ([]rawMapperInfo, error) {
	var mappers []rawMapperInfo //nolint:prealloc
	// For the ast declarations we care about (juryrig ones)...
	for _, decl := range astFile.Decls {
		if !isJuryRigMapperDecl(decl) {
			continue
		}

		// ...Extract some raw details we'll need from the declaration.
		raw, err := extractRawMapperInfo(fset, astFile, body, decl)
		if err != nil {
			return nil, fmt.Errorf("could not extract mapper: %w", err)
		}

		mappers = append(mappers, raw)
	}

	return mappers, nil
}

func isJuryRigMapperDecl(decl ast.Decl) bool {
	// A declaration is a juryrig declaration if it has
	// juryrig comments
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return false
	}

	return isJuryRigCommentGroup(genDecl.Doc)
}

func extractRawMapperInfo(fset *token.FileSet, astFile *ast.File, body []byte, decl ast.Decl) (rawMapperInfo, error) {
	// Extract details...
	genDecl, typeSpec, intSpec, err := extractRawMapperInfoAstDetails(fset, decl)
	if err != nil {
		return rawMapperInfo{}, err
	}

	fnInfos, err := extractRawMapperFuncInfos(fset, astFile, body, intSpec.Methods.List)

	if err != nil {
		return rawMapperInfo{}, fmt.Errorf("could not extract methods for mapper: %w", err)
	}

	comments := filterTaggedComments(genDecl.Doc.List, juryRigTag)

	// ...and Map.
	return rawMapperInfo{
		name:          typeSpec.Name.Name,
		topJrComments: comments,
		fns:           fnInfos,
	}, nil
}

func extractRawMapperInfoAstDetails(
	fset *token.FileSet,
	decl ast.Decl,
) (*ast.GenDecl, *ast.TypeSpec, *ast.InterfaceType, error) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil, nil, nil, fmt.Errorf("mapper tagged type is not a declaration %s: %w",
			locationDebugInfo(fset, decl), ErrUnexpectedAST)
	}

	if len(genDecl.Specs) != 1 {
		return nil, nil, nil, fmt.Errorf("expect mapper declaration to have 1 spec, but has %d: %w",
			len(genDecl.Specs), ErrUnexpectedAST)
	}

	spec := genDecl.Specs[0]
	typeSpec, ok := spec.(*ast.TypeSpec)

	if !ok {
		return nil, nil, nil, fmt.Errorf("mapper spec is not a type %s: %w",
			locationDebugInfo(fset, genDecl), ErrUnexpectedAST)
	}

	intSpec, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return nil, nil, nil, fmt.Errorf("mapper type is not an interface %s: %w",
			locationDebugInfo(fset, typeSpec), ErrSpec)
	}

	return genDecl, typeSpec, intSpec, nil
}

func extractRawMapperFuncInfos(
	fset *token.FileSet,
	astFile *ast.File,
	body []byte,
	methodFields []*ast.Field,
) ([]rawMapperFuncInfo, error) {
	// Loop through and delegate
	result := make([]rawMapperFuncInfo, len(methodFields))

	for i, methodField := range methodFields {
		fnInfo, err := extractRawMapperFuncInfo(fset, astFile, body, methodField)
		if err != nil {
			return nil, fmt.Errorf("could not extract method field %d: %w", i, err)
		}

		result[i] = fnInfo
	}

	return result, nil
}

func extractRawMapperFuncInfo(
	fset *token.FileSet,
	astFile *ast.File,
	body []byte,
	methodField *ast.Field,
) (rawMapperFuncInfo, error) {
	// Extract details...
	if len(methodField.Names) != 1 {
		return rawMapperFuncInfo{}, fmt.Errorf("expected method field to have 1 name, but has %d %s: %w",
			len(methodField.Names), locationDebugInfo(fset, methodField), ErrUnexpectedAST)
	}

	name := methodField.Names[0].Name
	funcType, ok := methodField.Type.(*ast.FuncType)

	if !ok {
		return rawMapperFuncInfo{}, fmt.Errorf("method field type is not a function %s: %w",
			locationDebugInfo(fset, methodField), ErrUnexpectedAST)
	}

	params, err := extractFuncParamaters(fset, astFile, body, funcType)

	if err != nil {
		return rawMapperFuncInfo{}, fmt.Errorf("could not extract func parameters: %w", err)
	}

	result, err := extractFuncResultType(astFile, body, funcType)

	if err != nil {
		return rawMapperFuncInfo{}, err
	}

	// ...and Map.
	return rawMapperFuncInfo{
		name:       name,
		parameters: params,
		result:     result,
		jrComments: filterTaggedComments(methodField.Doc.List, juryRigTag),
	}, nil
}

func extractFuncParamaters(
	fset *token.FileSet,
	astFile *ast.File,
	body []byte,
	fn *ast.FuncType,
) ([]Parameter, error) {
	result := make([]Parameter, len(fn.Params.List))
	// For each function parameter...
	for i, paramField := range fn.Params.List {
		// ...Read the name...
		if len(paramField.Names) != 1 {
			return nil, fmt.Errorf("paramater %d does not have one name %s: %w",
				i, locationDebugInfo(fset, paramField), ErrUnexpectedAST)
		}

		name := paramField.Names[0].Name

		// ...and Read the type.
		typ := readAsString(astFile, body, paramField.Type)

		result[i] = Parameter{
			Name: name,
			Type: typ,
		}
	}

	return result, nil
}

func extractFuncResultType(astFile *ast.File, body []byte, fn *ast.FuncType) (string, error) {
	if len(fn.Results.List) != 1 {
		return "", fmt.Errorf("mapper function must have exactly one result: %w",
			ErrSpec)
	}

	return readAsString(astFile, body, fn.Results.List[0].Type), nil
}

func isJuryRigCommentGroup(commentGroup *ast.CommentGroup) bool {
	// A comment group is a juryrig comment croup if any comment
	// is a juryrig comment.
	if commentGroup == nil {
		return false
	}

	for _, cmt := range commentGroup.List {
		if isTaggedComment(cmt, juryRigTag) {
			return true
		}
	}

	return false
}

func filterTaggedComments(comments []*ast.Comment, tag string) []string {
	var result []string

	for _, cmt := range comments {
		if isTaggedComment(cmt, tag) {
			result = append(result, strings.TrimSpace(cmt.Text))
		}
	}

	return result
}

func isTaggedComment(comment *ast.Comment, tag string) bool {
	return strings.HasPrefix(strings.TrimSpace(comment.Text), tag)
}

func locationDebugInfo(fset *token.FileSet, node ast.Node) string {
	position := fset.Position(node.Pos())
	return fmt.Sprintf("[file %s, line %d]", position.Filename, position.Line)
}

// Take an ast node and read the actual related source.
func readAsString(astFile *ast.File, body []byte, node ast.Node) string {
	offset := astFile.Pos()
	return string(body[node.Pos()-offset : node.End()-offset])
}
