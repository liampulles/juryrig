package parse

import (
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
	parameters map[string]string
	results    []string
	jrComments []string
}

// Just extract the most basic raw details from the files. Keep the
// ast stuff here basically.
func extractRaw(filename string) ([]rawMapperInfo, error) {
	// Read ast and files
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse file %s: %w", filename, err)
	}
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}

	// Parse
	result, err := parseMappers(fset, body, astFile)
	if err != nil {
		return nil, fmt.Errorf("could not parse ast: %w", err)
	}
	return result, nil
}

func parseMappers(fset *token.FileSet, body []byte, astFile *ast.File) ([]rawMapperInfo, error) {
	var result []rawMapperInfo
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
		result = append(result, raw)
	}
	return result, nil
}

func isJuryRigMapperDecl(decl ast.Decl) bool {
	// A declaration is a juryrig declaration if it has
	// juryrig comments
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return false
	}
	return isJuryRigMapperCommentGroup(genDecl.Doc)
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

func extractRawMapperInfoAstDetails(fset *token.FileSet, decl ast.Decl) (*ast.GenDecl, *ast.TypeSpec, *ast.InterfaceType, error) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil, nil, nil, fmt.Errorf("mapper tagged type is not a declaration %s",
			locationDebugInfo(fset, decl))
	}

	if len(genDecl.Specs) != 1 {
		return nil, nil, nil, fmt.Errorf("expect mapper declaration to have 1 spec, but has %d",
			len(genDecl.Specs))
	}
	spec := genDecl.Specs[0]
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, nil, nil, fmt.Errorf("mapper spec is not a type %s",
			locationDebugInfo(fset, genDecl))
	}

	intSpec, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return nil, nil, nil, fmt.Errorf("mapper type is not an interface %s",
			locationDebugInfo(fset, typeSpec))
	}

	return genDecl, typeSpec, intSpec, nil
}

func extractRawMapperFuncInfos(fset *token.FileSet, astFile *ast.File, body []byte, methodFields []*ast.Field) ([]rawMapperFuncInfo, error) {
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

func extractRawMapperFuncInfo(fset *token.FileSet, astFile *ast.File, body []byte, methodField *ast.Field) (rawMapperFuncInfo, error) {
	// Extract details...
	if len(methodField.Names) != 1 {
		return rawMapperFuncInfo{}, fmt.Errorf("expected method field to have 1 name, but has %d %s",
			len(methodField.Names), locationDebugInfo(fset, methodField))
	}
	name := methodField.Names[0].Name
	funcType, ok := methodField.Type.(*ast.FuncType)
	if !ok {
		return rawMapperFuncInfo{}, fmt.Errorf("method field type is not a function %s",
			locationDebugInfo(fset, methodField))
	}
	params, err := extractFuncParamaters(fset, astFile, body, funcType)
	if err != nil {
		return rawMapperFuncInfo{}, fmt.Errorf("could not extract func parameters: %w", err)
	}

	// ...and Map.
	return rawMapperFuncInfo{
		name:       name,
		parameters: params,
		results:    extractFuncResultTypes(fset, astFile, body, funcType),
		jrComments: filterTaggedComments(methodField.Doc.List, juryRigTag),
	}, nil
}

func extractFuncParamaters(fset *token.FileSet, astFile *ast.File, body []byte, fn *ast.FuncType) (map[string]string, error) {
	result := make(map[string]string)
	// For each function parameter...
	for i, paramField := range fn.Params.List {
		// ...Read the name...
		if len(paramField.Names) != 1 {
			return nil, fmt.Errorf("paramater %d does not have one name %s",
				i, locationDebugInfo(fset, paramField))
		}
		name := paramField.Names[0].Name
		// ...and Read the type.
		typ := readAsString(fset, astFile, body, paramField.Type)

		result[name] = typ
	}
	return result, nil
}

func extractFuncResultTypes(fset *token.FileSet, astFile *ast.File, body []byte, fn *ast.FuncType) []string {
	var types []string
	for _, resultField := range fn.Results.List {
		typ := readAsString(fset, astFile, body, resultField.Type)
		types = append(types, typ)
	}
	return types
}

func isJuryRigMapperCommentGroup(commentGroup *ast.CommentGroup) bool {
	// A comment group is a juryrig comment croup if any comment
	// is a juryrig comment.
	for _, cmt := range commentGroup.List {
		if isTaggedComment(cmt, juryRigMapperTag) {
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
func readAsString(fset *token.FileSet, astFile *ast.File, body []byte, node ast.Node) string {
	offset := astFile.Pos()
	return string(body[node.Pos()-offset : node.End()-offset])
}
