package wire

import "github.com/liampulles/juryrig/internal/config"

// Wiring holds references to all the "services" available to
// the runtime
type Wiring struct {
	configSvc config.Service
}