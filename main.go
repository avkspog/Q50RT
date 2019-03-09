package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type ServerConfig struct {
	Host          string
	TelemetryPort string
	APIPort       string
	Version       string
	LogFileName   string
}

type Starter struct {
	waitGroup              *sync.WaitGroup
	onStartTelemetryServer func(wg *sync.WaitGroup)
	onStartAPIService      func(wg *sync.WaitGroup)
}

var serverConfig *ServerConfig

var LocalCache *Cache

func init() {
	serverConfig = new(ServerConfig)
	serverConfig.Version = "0.0.1.12"
	serverConfig.LogFileName = "q50tlm.log"
	serverConfig.Host = "127.0.0.1"
	serverConfig.TelemetryPort = "30731"
	serverConfig.APIPort = "30732"

	flag.StringVar(&serverConfig.Host, "host", "127.0.0.1", "-host=127.0.0.1")
	flag.StringVar(&serverConfig.TelemetryPort, "tlm_port", "30731", "-tlm_port=30731")
	flag.StringVar(&serverConfig.APIPort, "api_port", "30732", "-api_port=30732")
	flag.Parse()
}

func main() {
	f, err := os.OpenFile(serverConfig.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
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
		onStartTelemetryServer: func(wg *sync.WaitGroup) {
			wg.Add(1)
			go startTelemetryServer(serverConfig, wg)
		},
		onStartAPIService: func(wg *sync.WaitGroup) {
			wg.Add(1)
			fmt.Println("api server test start")
			wg.Done()
		},
	}
	starter.run()
}

func (s *Starter) run() {
	s.onStartTelemetryServer(s.waitGroup)
	s.onStartAPIService(s.waitGroup)
	s.waitGroup.Wait()
}

func (c *ServerConfig) telemetryAddr() string {
	return c.Host + ":" + c.TelemetryPort
}

func (c *ServerConfig) APIAddr() string {
	return c.Host + ":" + c.APIPort
}
