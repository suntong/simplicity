// !!! !!!
// WARNING: Code automatically generated. Editing discouraged.
// !!! !!!

package main

import (
	"flag"
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const progname = "simplicity" // os.Args[0]

// The Options struct defines the structure to hold the commandline values
type Options struct {
	configFile string // path to config `file`
	maxProcs   int    // GOMAXPROCS, default is NumCpu()
	debug      int    // debugging `level`
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

// Opts holds the actual values from the command line paramters
var Opts Options

////////////////////////////////////////////////////////////////////////////
// Commandline definitions

func init() {

	// set default values for command line paramters
	flag.StringVar(&Opts.configFile, "Conf", "simplicity.conf",
		"path to config `file`")
	flag.IntVar(&Opts.maxProcs, "MaxProcs", -1,
		"GOMAXPROCS, default is NumCpu()")
	flag.IntVar(&Opts.debug, "debug", 0,
		"debugging `level`")

	// Now override those default values from environment variables
	if len(Opts.configFile) == 0 ||
		len(os.Getenv("SIMPLICITY_CONF")) != 0 {
		Opts.configFile = os.Getenv("SIMPLICITY_CONF")
	}

}

// The Usage function shows help on commandline usage
func Usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage:\n %s [flags ...] \n\nFlags:\n\n",
		progname)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		"\n")
	os.Exit(0)
}
