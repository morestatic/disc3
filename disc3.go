package main

import (
	"fmt"
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"

	"runtime/pprof"

	"deepsolutionsvn.com/disc/cli"
)

var options = new(cli.Options)

const (
	cpuProfile  = "disc3.cpuprof"
	heapProfile = "disc3.memprof"
)

func startProfiling() (func(), error) {
	// start CPU profiling as early as possible
	ofi, err := os.Create(cpuProfile)
	if err != nil {
		return nil, err
	}
	err = pprof.StartCPUProfile(ofi)
	if err != nil {
		ofi.Close()
		return nil, err
	}
	go func() {
		for range time.NewTicker(time.Second * 30).C {
			err := writeHeapProfileToFile()
			if err != nil {
				panic(err)
			}
		}
	}()

	stopProfiling := func() {
		pprof.StopCPUProfile()
		ofi.Close() // captured by the closure
	}
	return stopProfiling, nil
}

func writeHeapProfileToFile() error {
	mprof, err := os.Create(heapProfile)
	if err != nil {
		return err
	}
	defer mprof.Close()
	return pprof.WriteHeapProfile(mprof)
}

func profileIfEnabled() (func(), error) {
	if false {
		stopProfilingFunc, err := startProfiling()
		if err != nil {
			return nil, err
		}
		return stopProfilingFunc, nil
	}
	return func() {}, nil
}

func main() {

	stop, err := profileIfEnabled()
	if err != nil {
		log.Fatal("Error while checking/starting profiling")
	}
	defer stop()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	p := flags.NewParser(options, flags.HelpFlag|flags.PassDoubleDash)

	// parse and execute commands
	_, err = p.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
