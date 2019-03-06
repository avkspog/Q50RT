package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	Version     = "0.0.1.12"
	LogFileName = "q50tlm.log"
)

var (
	host          = "127.0.0.1"
	telemetryPort = "8002"
	gRPCPort      = "9002"
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "-host=127.0.0.1")
	flag.StringVar(&telemetryPort, "tlm_port", "8002", "-tlm_port=8002")
	flag.StringVar(&gRPCPort, "grpc_port", "9002", "-grpc_port=9002")
	flag.Parse()
}

var LocalCache *Cache

type Starter struct {
	waitGroup              *sync.WaitGroup
	onStartTelemetryServer func(addr string, wg *sync.WaitGroup)
	onStartApiService      func(addr string, wg *sync.WaitGroup)
}

func main() {

	f, err := os.OpenFile(LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	LocalCache = NewCache()

	starter := &Starter{
		waitGroup: &sync.WaitGroup{},
		onStartTelemetryServer: func(addr string, wg *sync.WaitGroup) {
			go startTelemetryServer(addr, wg)
		},
		onStartApiService: func(addr string, wg *sync.WaitGroup) {
			wg.Done()
		},
	}
	starter.run()
}

func (s *Starter) run() {
	tlmAddr := host + ":" + telemetryPort
	apiAddr := host + ":" + gRPCPort

	s.waitGroup.Add(2)
	s.onStartTelemetryServer(tlmAddr, s.waitGroup)
	s.onStartApiService(apiAddr, s.waitGroup)
	s.waitGroup.Wait()
}
