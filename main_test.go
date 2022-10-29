package main_test

import (
	"os"
	"testing"

	goConfig "github.com/liampulles/go-config"
	"github.com/liampulles/juryrig/internal/wire"
	"github.com/stretchr/testify/assert"
)

func TestJuryrig_ValidExample(t *testing.T) {
	// Setup fixture
	args := []string{"juryrig", "gen", "-o", "actual.go"}
	cfgSource := goConfig.MapSource{
		"GOFILE": "testdata/film/mapper.go",
	}

	// Setup expectations
	expected, err := os.ReadFile("testdata/film/expected.go")
	assert.NoError(t, err, "could not read expected")

	// Exercise SUT
	code := wire.Run(args, cfgSource)

	// Verify results
	assert.Equal(t, 0, code, "non-zero exit code")
	actual, err := os.ReadFile("testdata/film/actual.go")
	assert.NoError(t, err, "could not read actual")
	assert.Equal(t, string(expected), string(actual))
}
