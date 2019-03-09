package main

import (
	pb "Q50RT/api"
	"context"
)

type APIPointServer struct {
}

func (s APIPointServer) Ping(ctx context.Context, ping *pb.PingCommand) (*pb.PingCommand, error) {
	return &pb.PingCommand{Message: "pong"}, nil
}

func newApiServer() *APIPointServer {
	s := &APIPointServer{}
	return s
}
