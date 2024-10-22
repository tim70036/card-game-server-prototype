package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
)

func main() {
	flag.CommandLine.SetOutput(io.Writer(os.Stdout))
	// flag.PrintDefaults()

	var runAsClient = flag.Bool("client", false, "run as client")
	flag.Parse()

	if *runAsClient {
		client, err := BuildClient()
		if err != nil {
			log.Fatalf("wire build client failed %v", err)
			return
		}
		client.Run()
	} else {
		go runProfiler()
		gameMaker, err := BuildGameMaker()
		if err != nil {
			log.Fatalf("wire build game maker failed %v", err)
			return
		}
		gameMaker.Run()
	}
}

func runProfiler() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	if err := http.ListenAndServe("localhost:9487", mux); err != nil {
		log.Fatalf("profiler ListenAndServe failed %v", err)
	}
}
