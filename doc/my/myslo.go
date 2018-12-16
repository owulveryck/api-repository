// Package my is an instanciation of the concept exposed in the talks
package my

import (
	"time"

	"github.com/owulveryck/api-repository/injector"
)

// SLO as defined before
// START_SLO OMIT
var SLO = &injector.SLO{
	Latency: map[float64]time.Duration{
		0.95: 400 * time.Millisecond,
		0.99: 750 * time.Millisecond,
	},
	Allowed5xxErrors: 100.0 - 97.0,
	Verbose:          true,
}

// END_SLO OMIT
