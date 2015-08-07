////////////////////////////////////////////////////////////////////////////
// Porgram: config.go
// Purpose: Go application configuration solution
// Authors: Tong Sun (c) 2015, All rights reserved
////////////////////////////////////////////////////////////////////////////

/*

Application configuration that satisfies:

- Have a set of default values defined in the config file
- Environment variables can override them
- Variables passed from the command line takes the highest priority

*/

package main

import (
	"flag"
	"runtime"

	"github.com/koding/multiconfig"
)

var (
	configFile = flag.String("conf", "simplicity.toml", "path to config file")
)

type Config struct {
	MaxLen   int
	MaxProcs int

	// [Webapp]
	Webapp struct {
		Path       string
		PathAdmin  string
		Credential string
	}
}

// ConfigGet will get configuration from multiple layers, toml file, environment variables, or command line.
func ConfigGet() *Config {
	// Start with an empty struct of the configuration
	cf := new(Config)
	m := multiconfig.NewWithPath(*configFile)
	// Populated the configuration struct
	m.MustLoad(cf) // Check for error
	//log.Printf("%+v\n", cf)
	//log.Printf("P: %d, B: '%s', F: '%s'\n", cf.MaxProcs, cf.Webapp.Path, *configFile)

	cf.ConfigApply()
	return cf
}

// ConfigApply will apply the configuration
func (cf *Config) ConfigApply() {
	if cf.MaxProcs <= 0 {
		cf.MaxProcs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(cf.MaxProcs)
}
