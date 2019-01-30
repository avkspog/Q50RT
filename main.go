package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	Version     = "0.0.1.11"
	LogFileName = "q50tlm.log"
)

var (
	Host          string = "127.0.0.1"
	TelemetryPort string = "8002"
	GRPCPort      string = "9002"
)

func init() {
	flag.StringVar(&Host, "host", "127.0.0.1", "-host=127.0.0.1")
	flag.StringVar(&TelemetryPort, "tlm_port", "8002", "-tlm_port=8002")
	flag.StringVar(&GRPCPort, "grpc_port", "9002", "-grpc_port=9002")
	flag.Parse()
}

var LocalCache *Cache

func main() {
	f, err := os.OpenFile(LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)

	log.SetOutput(mw)

	LocalCache = NewCache()

	StartTelemetryServer(Host + ":" + TelemetryPort)
}
