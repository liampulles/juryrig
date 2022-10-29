package template

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/liampulles/juryrig/internal/parse"
)

//nolint:gochecknoglobals
var mapperFileTemplate = template.Must(template.New("mapper-file").
	Parse(`package {{ .Package }}{{ range $m, $mapper := .Mappers }}
type {{$mapper.Name}}Impl struct{}{{ range $f, $func := $mapper.Functions }}
func (impl *{{$mapper.Name}}Impl) {{$func.Name}}({{$func.Params}}) {{$func.Result}} {
	return {{$func.Result}}{
		{{ range $func.Directives }}{{.}},
		{{ end }}
	}
}
{{ end }}{{ end }}`))

func Generate(in parse.JuryrigSpec) []byte {
	spec := mapSpec(in)
	w := &bytes.Buffer{}

	// Template it
	if err := mapperFileTemplate.Execute(w, spec); err != nil {
		return []byte(fmt.Sprintf("<<<TEMPLATE ERROR: %s>>>", err.Error()))
	}

	// Format it
	formatted, err := format.Source(w.Bytes())
	if err != nil {
		return []byte(fmt.Sprintf("<<<FORMAT ERROR: %s>>>", err.Error()))
	}

	return formatted
}

func mapSpec(in parse.JuryrigSpec) spec {
	mappers := make([]mapper, len(in.Mappers))
	for i, mapper := range in.Mappers {
		mappers[i] = mapMapper(mapper)
	}

	return spec{
		Package: in.Package,
		Mappers: mappers,
	}
}

func mapMapper(in parse.Mapper) mapper {
	funcs := make([]function, len(in.MapperFunctions))
	for i, fn := range in.MapperFunctions {
		funcs[i] = mapFunc(fn)
	}

	return mapper{
		Name:      in.Name,
		Functions: funcs,
	}
}

func mapFunc(in parse.MapperFunction) function {
	params := make([]string, len(in.Function.Parameters))
	for i, param := range in.Function.Parameters {
		params[i] = mapParam(param)
	}

	directives := make([]string, len(in.Directives))
	for i, directive := range in.Directives {
		directives[i] = mapDirective(directive)
	}

	return function{
		Name:       in.Function.Name,
		Params:     strings.Join(params, ", "),
		Result:     in.Function.Result,
		Directives: directives,
	}
}

func mapParam(in parse.Parameter) string {
	return fmt.Sprintf("%s %s", in.Name, in.Type)
}

func mapDirective(in parse.Directive) string {
	switch v := in.(type) {
	case parse.LinkDirective:
		return formatDirective(v.Target, mapSource(v.Source))
	case parse.LinkFuncDirective:
		return formatDirective(v.Target, mapLinkFuncValue(v))
	case parse.IgnoreDirective:
		return fmt.Sprintf("// %s", formatDirective(v.Target, "(ignored)"))
	}
	// Should be handled by parse stage...
	return fmt.Sprintf("<<ERROR: UNKNOWN DIRECTIVE TYPE %T>>", in)
}

func formatDirective(target parse.Target, value string) string {
	return fmt.Sprintf("%s: %s", target.Field, value)
}

func mapLinkFuncValue(in parse.LinkFuncDirective) string {
	sources := make([]string, len(in.Sources))
	for i, source := range in.Sources {
		sources[i] = mapSource(source)
	}

	return fmt.Sprintf("impl.%s(%s)", in.FunctionName, strings.Join(sources, ", "))
}

func mapSource(in parse.Source) string {
	if in.Field != "" {
		return fmt.Sprintf("%s.%s", in.Parameter, in.Field)
	}

	return in.Parameter
}

type spec struct {
	Package string
	Mappers []mapper
}

type mapper struct {
	Name      string
	Functions []function
}

type function struct {
	Name       string
	Params     string
	Result     string
	Directives []string
}
