package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/avkspog/brts"
)

var (
	host string
	port string
)

func main() {
	tcpServer := brts.Create(getTCPAddress())
	tcpServer.SetTimeout(3 * time.Minute)
	tcpServer.SetMessageDelim(']')

	tcpServer.OnServerStarted(func(addr *net.TCPAddr) {
		log.Printf("Q50Watch telemetry server started on address: %v", addr.String())
	})

	tcpServer.OnServerStopped(func() {
		log.Println("Q50Watch telemetry server stopped")
	})

	tcpServer.OnNewConnection(func(c *brts.Client) {
		log.Printf("accepted connection from: %v", c.Conn.RemoteAddr())
	})

	tcpServer.OnMessageReceive(func(c *brts.Client, data *[]byte) {
		
	})

	tcpServer.OnConnectionLost(func(c *brts.Client) {
		log.Printf("closing connection from %v", c.Conn.RemoteAddr())
	})

	if err := tcpServer.Start(); err != nil {
		log.Printf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "-host=127.0.0.1")
	flag.StringVar(&port, "port", "8002", "-port=8002")
	flag.Parse()
}

func getTCPAddress() string {
	return host + ":" + port
}
