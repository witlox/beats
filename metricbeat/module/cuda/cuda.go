package cuda

import (
	"github.com/elastic/beats/metricbeat/mb"
)

func init() {
	// Register the ModuleFactory function for the "cuda" module.
	if err := mb.Registry.AddModule("cuda", mb.NewModule); err != nil {
		panic(err)
	}
}
