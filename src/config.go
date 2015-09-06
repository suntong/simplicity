////////////////////////////////////////////////////////////////////////////
// Porgram: config.go
// Purpose: Go application configuration
// Authors: Tong Sun (c) 2015, All rights reserved
////////////////////////////////////////////////////////////////////////////

/*

Based on https://github.com/pkieltyka/godo-app/blob/master/config.go
by Peter Kieltyka.

*/

package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
)

var ErrNoConfigFile = errors.New("No configuration file specified.")

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

func ConfigGet() *Config {
	var err error
	var cf *Config = NewConfig()

	flag.Usage = Usage
	flag.Parse()

	// Load config file, from flag or env (if specified)
	_, err = cf.ConfigFromFile(Opts.configFile, os.Getenv("APPCONFIG"))
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("P: %d, B: '%s', F: '%s'\n", cf.MaxProcs, cf.Webapp.Path)

	cf.ConfigApply()

	return cf
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigFromFile(confFile string, confEnv string) (*Config, error) {
	cf := &Config{}
	return cf.ConfigFromFile(confFile, confEnv)
}

func (cf *Config) ConfigFromFile(confFile string, confEnv string) (*Config, error) {
	if confFile == "" {
		confFile = confEnv
	}

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		return nil, ErrNoConfigFile
	}

	if _, err := toml.DecodeFile(confFile, &cf); err != nil {
		return nil, err
	}
	return cf, nil
}

func (cf *Config) ConfigApply() {
	if cf.MaxProcs <= 0 {
		cf.MaxProcs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(cf.MaxProcs)
}
