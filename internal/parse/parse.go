package parse

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Read parses Go source into Mapper definitions.
func Read(filename string) (JuryrigSpec, error) {
	// Just read the raw details (don't want to deal with ast stuff here)
	pkg, raw, err := extractRaw(filename)
	if err != nil {
		return JuryrigSpec{}, fmt.Errorf("could not extract raw: %w", err)
	}

	// Convert to Mapper types
	mappers := make([]Mapper, len(raw))

	for i, rawI := range raw {
		mapper, err := convertRawToMapper(rawI)
		if err != nil {
			return JuryrigSpec{}, err
		}

		mappers[i] = mapper
	}

	// Join
	return JuryrigSpec{
		Package: pkg,
		Mappers: mappers,
	}, nil
}

func convertRawToMapper(raw rawMapperInfo) (Mapper, error) {
	// Map pieces...
	mapperFuncs, err := convertRawFuncsToMapperFuncs(raw.fns)
	if err != nil {
		return Mapper{}, fmt.Errorf("cannot create mapper for %s: %w",
			raw.name, err)
	}

	// ...and join.
	return Mapper{
		Name:            raw.name,
		MapperFunctions: mapperFuncs,
	}, nil
}

func convertRawFuncsToMapperFuncs(rawFuncs []rawMapperFuncInfo) ([]MapperFunction, error) {
	mapperFuncs := make([]MapperFunction, len(rawFuncs))

	for i, rawFunc := range rawFuncs {
		mapperFunc, err := convertRawFuncToMapperFunc(rawFunc)
		if err != nil {
			return nil, err
		}

		mapperFuncs[i] = mapperFunc
	}

	return mapperFuncs, nil
}

func convertRawFuncToMapperFunc(rawFunc rawMapperFuncInfo) (MapperFunction, error) {
	// Map pieces...
	mapperFn := createMapperFunction(rawFunc)
	directives, err := createDirectives(rawFunc.jrComments)

	if err != nil {
		return MapperFunction{}, err
	}

	// ...and join.
	return MapperFunction{
		Function:   mapperFn,
		Directives: directives,
	}, nil
}

func createMapperFunction(rawFunc rawMapperFuncInfo) Function {
	return Function{
		Name:       rawFunc.name,
		Parameters: rawFunc.parameters,
		Result:     rawFunc.result,
	}
}

func createDirectives(jrComments []string) ([]Directive, error) {
	// Each comment should correspond to one directive.
	directives := make([]Directive, len(jrComments))

	for i, jrComment := range jrComments {
		directive, err := createDirective(jrComment)
		if err != nil {
			return nil, err
		}

		directives[i] = directive
	}

	return directives, nil
}

// Example: `+juryrig:link:ef.runtime->runtime`.
var juryrigDirectiveRegex = regexp.MustCompile(`\+juryrig:(\w+):(.+)`)

func createDirective(jrComment string) (Directive, error) {
	// Parse the comment for raw details
	var name, details string
	if err := extractRegex(juryrigDirectiveRegex, jrComment, &name, &details); err != nil {
		return nil, fmt.Errorf("[%s] is not a valid juryrig directive: %w",
			jrComment, ErrSpec)
	}

	// Delegate to more specific directive parsing
	switch name {
	case "link":
		return createLinkDirective(details)
	case "linkfunc":
		return createLinkFuncDirective(details)
	case "ignore":
		return createIgnoreDirective(details)
	default:
		return nil, fmt.Errorf("[%s] does not contain a recognized directive: %w",
			jrComment, ErrSpec)
	}
}

var juryrigLinkDetailsRegex = regexp.MustCompile(`^(.+)->(\w+)$`)

func createLinkDirective(details string) (LinkDirective, error) {
	// Parse the details...
	var sourceStr, targetStr string
	if err := extractRegex(juryrigLinkDetailsRegex, details, &sourceStr, &targetStr); err != nil {
		return LinkDirective{}, fmt.Errorf("[%s] is not valid config for the link directive: %w",
			details, ErrSpec)
	}

	source := parseSource(sourceStr)

	// ...and join
	return LinkDirective{
		Source: source,
		Target: Target{
			Field: targetStr,
		},
	}, nil
}

var juryrigLinkFuncDetailsRegex = regexp.MustCompile(`^(\w+)->(\w+)->(\w+)$`)

func createLinkFuncDirective(details string) (LinkFuncDirective, error) {
	// Parse the details...
	var from, fn, target string
	if err := extractRegex(juryrigLinkFuncDetailsRegex, details, &from, &fn, &target); err != nil {
		return LinkFuncDirective{}, fmt.Errorf("[%s] is not valid config for the linkfunc directive: %w",
			details, ErrSpec)
	}
	// (can be multiple parameters, comma separated)
	sourceStrsRaw := strings.Split(from, ",")
	sourceStrs := mapStrings(sourceStrsRaw, strings.TrimSpace)
	sources := make([]Source, len(sourceStrs))

	for i, sourceStr := range sourceStrs {
		sources[i] = parseSource(sourceStr)
	}

	// ...and join.
	return LinkFuncDirective{
		Sources:      sources,
		FunctionName: fn,
		Target: Target{
			Field: target,
		},
	}, nil
}

func createIgnoreDirective(details string) (IgnoreDirective, error) {
	// The details in this case should just be a target.
	if len(details) == 0 {
		return IgnoreDirective{}, fmt.Errorf("[%s] is not valid config for the ignore directive: %w",
			details, ErrSpec)
	}

	return IgnoreDirective{
		Target: Target{
			Field: details,
		},
	}, nil
}

var juryrigFieldSourceRegex = regexp.MustCompile(`^(\w+).(\w+)$`)

func parseSource(str string) Source {
	// Try it as a paramater-field variant
	var parameter, field string
	err := extractRegex(juryrigFieldSourceRegex, str, &parameter, &field)

	if err == nil {
		return Source{
			Parameter: parameter,
			Field:     field,
		}
	}

	// Okay, assume it is a parameter only variant.
	return Source{
		Parameter: str,
		Field:     "",
	}
}

// >>> String/regex helpers <<<

func mapStrings(in []string, fn func(string) string) []string {
	out := make([]string, len(in))
	for i, str := range in {
		out[i] = fn(str)
	}

	return out
}

var ErrRegexMismatch = errors.New("regex mismatch")

func extractRegex(r *regexp.Regexp, str string, into ...*string) error {
	matches := r.FindStringSubmatch(str)
	if len(matches) < len(into)+1 {
		return ErrRegexMismatch
	}

	for i, intoI := range into {
		*intoI = matches[i+1]
	}

	return nil
}
