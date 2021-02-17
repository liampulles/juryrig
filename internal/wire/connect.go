package wire

import (
	goConfig "github.com/liampulles/go-config"
	"github.com/liampulles/juryrig/internal/config"
)

// Connect wires up services to produce Wiring, which can be
// injected as required
func Connect(source goConfig.Source) *Wiring {
	configSvc := config.NewServiceImpl(source)

	return &Wiring{
		configSvc: configSvc,
	}
}
