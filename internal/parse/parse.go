package parse

import (
	"fmt"
)

const (
	juryRigTag       = "// +juryrig:"
	juryRigMapperTag = juryRigTag + "mapper:"
)

// ParseFileWithMapperDefs parses Go source into Mapper definitions
func ParseFileWithMapperDefs(filename string) ([]Mapper, error) {
	// Just read the raw details
	raw, err := extractRaw(filename)
	if err != nil {
		return nil, fmt.Errorf("could not extract raw: %w", err)
	}

	// Convert to Mapper types
	result := make([]Mapper, len(raw))
	for i, rawI := range raw {
		result[i] = convertRawToMapper(rawI)
	}
	return result, nil
}

func convertRawToMapper(raw rawMapperInfo) Mapper {
	// TODO:
	return Mapper{}
}
