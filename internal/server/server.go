package main

import (
	"context"
	"net"
	"os"
	"strconv"

	"fmt"
	"log"

	"groove-go.nyu.edu/pkg/crypt"
	parser "groove-go.nyu.edu/pkg/parser"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "groove-go.nyu.edu/route"
)

type clientKey struct {
	clientId     int
	symmetricKey []byte
}

type server struct {
	pb.UnimplementedServerServer
	s             *parser.Server
	publicKey     crypt.PublicKey
	privateKey    crypt.PrivateKey
	circuitKeyMap map[int][]byte
}

func (s *server) Initialize() error {
	var err error
	s.publicKey, s.privateKey, err = crypt.GenerateKeys()
	if err != nil {
		return err
	}
	s.circuitKeyMap = make(map[int][]byte)

	return nil
}

func (s *server) CircuitSetup(ctx context.Context, req *pb.CircuitSetupRequest) (*pb.CircuitSetupResponse, error) {
	csl := parser.ReadCircuitSetupLayer(req.Message)
	symmetricKey := crypt.ComputeSymmetricKey(crypt.UnmarshalPublicKey(csl.EphPublicKey), s.privateKey)
	cfd, err := crypt.DecryptCircuitSetupLayer(csl.Data, symmetricKey)
	if err != nil {
		return nil, err
	}

	var res *pb.CircuitSetupResponse
	var returnData []byte

	if cfd.NextHopPort != 0 {
		nextHopServer, err := parser.FetchServerByPort(cfd.NextHopPort)
		if err != nil {
			return nil, err
		}

		// Dial next hop server
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", nextHopServer.Host, nextHopServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to dial: %v", err)
		}
		defer conn.Close()

		nhs := pb.NewServerClient(conn)

		res, err = nhs.CircuitSetup(ctx, &pb.CircuitSetupRequest{Message: cfd.Data})
		if err != nil {
			return nil, err
		}
		returnData = res.Message
	} else {
		fmt.Printf("Dead drop received at endpoint: %s\n", string(cfd.Data))
		returnData = cfd.Data // should contain dead drop ID
	}

	returnMsg, err := crypt.EncryptSymmetric(returnData, symmetricKey)
	if err != nil {
		return nil, err
	}

	// Save tag if success
	s.addCircuitKey(cfd.Tag, symmetricKey)

	res = &pb.CircuitSetupResponse{Message: returnMsg}

	return res, nil
}

func (s *server) FetchPublicKey(ctx context.Context, req *pb.FetchPublicKeyRequest) (*pb.FetchPublicKeyResponse, error) {
	// TODO: handle empty public key?
	return &pb.FetchPublicKeyResponse{PublicKey: crypt.MarshalPublicKey(s.publicKey)}, nil
}

func (s *server) addCircuitKey(tag int, symmetricKey []byte) {
	s.circuitKeyMap[tag] = symmetricKey
}

func main() {
	// Setup log
	f, err := os.OpenFile("log.out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	server := server{}
	server.Initialize()

	// Read port from command line
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		println("Could not parse port number: %v", err)
	}

	server.s, err = parser.FetchServerByPort(port)
	if err != nil {
		log.Fatalf("Failed to fetch server with port %d: %v", port, err)
	}

	// Start server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.s.Host, server.s.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	svr := grpc.NewServer()
	pb.RegisterServerServer(svr, &server)
	log.Printf("Server %d listening at %v", server.s.Id, lis.Addr())

	svr.Serve(lis)
}
