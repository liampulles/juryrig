package config

// Config defines environment configuration given to the app.
type Config struct {
	Architecture    string
	OperatingSystem string
	BaseFilename    string
	DirectiveLine   int
	PackageName     string
}
