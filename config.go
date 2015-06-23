package rex

import (
	"flag"
	"runtime"
)

var (
	debug    bool
	port     int
	maxprocs int
)

func configure() {
	flag.BoolVar(&debug, "debug", Env.Bool("DEBUG", true), "flag to toggle debug mode")
	flag.IntVar(&port, "port", Env.Int("PORT", 5000), "port to run the application server")
	flag.IntVar(&maxprocs, "maxprocs", Env.Int("MAXPROCS", runtime.NumCPU()), "maximum cpu processes to run the server")
	flag.Parse()
}
