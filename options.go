package rex

import (
	"flag"
	"runtime"
	"sync"

	"github.com/goanywhere/env"
)

var (
	debug    bool
	port     int
	maxprocs int
	once     sync.Once
)

func configure() {
	once.Do(func() {
		flag.BoolVar(&debug, "debug", env.Bool("DEBUG", true), "flag to toggle debug mode")
		flag.IntVar(&port, "port", env.Int("PORT", 5000), "port to run the application server")
		flag.IntVar(&maxprocs, "maxprocs", env.Int("MAXPROCS", runtime.NumCPU()), "maximum cpu processes to run the server")
		flag.Parse()
	})
}
