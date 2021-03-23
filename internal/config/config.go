package config

import (
	goConfig "github.com/liampulles/go-config"
)

// Config defines environment configuration given to the app.
type Config struct {
	Architecture    string
	OperatingSystem string
	BaseFilename    string
	DirectiveLine   int
	PackageName     string
}

// Service allows for interfaceng with environment config
type Service interface {
	Read() (*Config, error)
}

// ServiceImpl implements service
type ServiceImpl struct {
	source goConfig.Source
}

// NewServiceImpl is a constructor
func NewServiceImpl(source goConfig.Source) *ServiceImpl {
	return &ServiceImpl{
		source: source,
	}
}

// Check we implement the interface
var _ Service = &ServiceImpl{}

// Read reads the environment into Config. If there are any issues,
// an error is returned.
func (s *ServiceImpl) Read() (*Config, error) {
	typedSource := goConfig.NewTypedSource(s.source)
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
