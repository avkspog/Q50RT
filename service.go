package main

import (
	pb "Q50RT/api"
	ps "Q50RT/q50"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
)

type APIServer struct {
	protocolVersion string
	address         string
	server          *grpc.Server
}

func (s *APIServer) Ping(ctx context.Context, ping *pb.PingCommand) (*pb.PingCommand, error) {
	return &pb.PingCommand{Message: s.protocolVersion}, nil
}

func (s *APIServer) LastPoint(ctx context.Context, idn *pb.Identifier) (*pb.Point, error) {
	if idn == nil {
		return &pb.Point{}, errors.New("Empty client identifier")
	}

	if len(idn.ClientId) == 0 {
		return &pb.Point{}, errors.New("Invalid client id")
	}

	if s.protocolVersion != idn.Version {
		return &pb.Point{}, fmt.Errorf("Protocol version %s not support", idn.Version)
	}

	msg, ok := LocalCache.Get(idn.ClientId)
	if !ok {
		return &pb.Point{}, nil
	}

	message, ok := msg.(ps.Message)
	if !ok {
		return &pb.Point{}, nil
	}

	point := &pb.Point{
		Version:        s.protocolVersion,
		MessageType:    message.MessageType,
		NetType:        message.NetType,
		DeviceId:       message.ID,
		BatteryPercent: uint32(message.BatteryPercent),
		ReceiveTime:    message.ReceiveTime.UnixNano(),
		DeviceTime:     message.DeviceTime.UnixNano(),
		Latitude:       message.Latitude,
		Longitude:      message.Longitude,
	}

	return point, nil
}

func (s *APIServer) ServerStatistic(ctx context.Context, command *pb.ServerCommand) (*pb.ServerResponse, error) {
	return nil, nil
}

func createAPIServer(c *ServerConfig) *APIServer {
	apis := &APIServer{
		protocolVersion: serverConfig.ProtocolVersion,
		address:         c.APIAddr(),
		server:          grpc.NewServer(),
	}

	pb.RegisterRoutePointServer(apis.server, apis)

	return apis
}

func StartAPIServer(c *ServerConfig, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		log.Println("Q50Watch api server stopped")
	}()

	s := createAPIServer(c)

	sigs := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
		<-sigs
		s.Shutdown()
	}()

	lis, err := net.Listen("tcp", c.APIAddr())
	if err != nil {
		log.Printf("Fatal error: %s", err.Error())
	}

	log.Printf("Q50Watch api server v%s started on address: %v", serverConfig.Version, c.APIAddr())
	if err := s.server.Serve(lis); err != nil {
		log.Printf("Fatal error: %s", err.Error())
	}
}

func (s *APIServer) Shutdown() {
	s.server.Stop()
}
