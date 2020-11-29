package config

import (
	goConfig "github.com/liampulles/go-config"
)

// Read reads the environment into Config. If there are any issues,
// an error is returned.
func Read(source goConfig.Source) (*Config, error) {
	typedSource := goConfig.NewTypedSource(source)
	config := &Config{}

	if err := goConfig.LoadProperties(typedSource,
		goConfig.StrProp("GOARCH", &config.Architecture, false),
		goConfig.StrProp("GOOS", &config.OperatingSystem, false),
		goConfig.StrProp("GOFILE", &config.BaseFilename, false),
		goConfig.IntProp("GOLINE", &config.DirectiveLine, false),
		goConfig.StrProp("GOPACKAGE", &config.PackageName, false),
	); err != nil {
		return nil, err
	}

	return config, nil
}
