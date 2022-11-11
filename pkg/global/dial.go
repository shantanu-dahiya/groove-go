package global

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "groove-go.nyu.edu/route"
)

func (c *Client) Dial() (pb.ClientClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", c.Host, c.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	return pb.NewClientClient(conn), conn
}

func (s *Server) Dial() (pb.ServerClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.Host, s.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	return pb.NewServerClient(conn), conn
}
