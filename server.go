package main

import (
	"Q50RT/q50"
	"fmt"
	"github.com/avkspog/brts"
	"log"
	"net"
	"os"
	"time"
)

func StartTelemetryServer(addr string) {
	tcpServer := brts.Create(addr)
	tcpServer.SetTimeout(3 * time.Minute)
	tcpServer.SetMessageDelim(']')

	tcpServer.OnServerStarted(func(addr *net.TCPAddr) {
		log.Printf("Q50Watch telemetry server v%s started on address: %v", Version, addr.String())
	})

	tcpServer.OnServerStopped(func() {
		log.Println("Q50Watch telemetry server stopped")
	})

	tcpServer.OnNewConnection(func(c *brts.Client) {
		log.Printf("accepted connection from: %v", c.Conn.RemoteAddr())
	})

	tcpServer.OnMessageReceive(func(c *brts.Client, data *[]byte) {
		s := fmt.Sprintf("%s", *data)
		log.Println(s)
		go process(data)
	})

	tcpServer.OnConnectionLost(func(c *brts.Client) {
		log.Printf("closing connection from %v", c.Conn.RemoteAddr())
	})

	if err := tcpServer.Start(); err != nil {
		log.Printf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func process(data *[]byte) {
	message, err := q50.Parse(data)
	if err != nil {
		log.Println(err)
		return
	}

	if message == nil {
		log.Println("message is nil")
		return
	}

	if message.ID == 0 {
		log.Println("message id is zero")
		return
	}

	cmsg, ok := LocalCache.Get(message.ID)
	if ok {
		cachedMessage := cmsg.(*q50.Message)

		//if cachedMessage.DeviceTime.Before(message.DeviceTime) {
		//	log.Printf("message device time before cached message time: %s", message.DeviceTime)
		//	return
		//}

		cachedMessage.MessageType = message.MessageType
		cachedMessage.NetType = message.NetType
		cachedMessage.ReceiveTime = message.ReceiveTime
		cachedMessage.DeviceTime = message.DeviceTime
		if message.BatteryPercent != 0 {
			cachedMessage.BatteryPercent = message.BatteryPercent
		}
		if message.Latitude != 0 && message.Longitude != 0 {
			cachedMessage.Latitude = message.Latitude
			cachedMessage.Longitude = message.Longitude
		}
		LocalCache.Set(message.ID, cachedMessage)
	} else {
		LocalCache.Set(message.ID, message)
	}

	log.Println(len(LocalCache.Items))
	for k, v := range LocalCache.Items {
		msg := v.Value.(*q50.Message)
		log.Printf("element %s - %v, id = %v, bat = %v, lat = %v, lon = %v", msg.MessageType,
			k, msg.ID, msg.BatteryPercent, msg.Latitude, msg.Longitude)
	}
}