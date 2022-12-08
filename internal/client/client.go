package main

import (
	"context"

	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"groove-go.nyu.edu/pkg/crypt"
	parser "groove-go.nyu.edu/pkg/parser"

	"google.golang.org/grpc"

	pb "groove-go.nyu.edu/route"
)

type client struct {
	pb.UnimplementedClientServer
	publicKey  crypt.PublicKey
	privateKey crypt.PrivateKey
	buddies    []*parser.Buddy
	c          *parser.Client
}

func (c *client) Initialize() error {
	var err error
	c.publicKey, c.privateKey, err = crypt.GenerateKeys()
	return err
}

func (c *client) SymmetricKeyGen(ctx context.Context, req *pb.SymmetricKeyGenRequest) (*pb.SymmetricKeyGenResponse, error) {
	peerPublicKey := crypt.UnmarshalPublicKey(req.PublicKey)
	c.addBuddy(int(req.ClientId), peerPublicKey)

	return &pb.SymmetricKeyGenResponse{PublicKey: crypt.MarshalPublicKey(c.publicKey)}, nil
}

func (c *client) addBuddy(clientId int, publicKey crypt.PublicKey) {
	symmetricKey := crypt.ComputeSymmetricKey(publicKey, c.privateKey)

	newBuddy := parser.NewBuddy(clientId, symmetricKey)
	for _, buddy := range c.buddies {
		if buddy.Id == clientId {
			// already added
			log.Printf("Client %d already added as buddy", clientId)
			return
		}
	}
	// If new
	c.buddies = append(c.buddies, newBuddy)
}

func (c *client) setUpCircuits(ctx context.Context) {
	var firstHop *parser.CircuitElement

	for _, b := range c.buddies {
		for _, circuit := range b.Circuits {
			log.Printf("Working on circuit %s for buddy %d of client %d", circuit.Print(), b.Id, c.c.Id)
			firstHop = (*circuit)[0]

			// Fetch public keys and derive symmetric keys for each server
			for _, e := range *circuit {
				server, serverConn := e.Server.Dial()

				pkRes, err := server.FetchPublicKey(ctx, &pb.FetchPublicKeyRequest{})
				if err != nil {
					log.Fatalf("Failed to fetch public key from server %d: %v", e.Server.Id, err)
				}

				// TODO: change to ephemeral
				e.SymmetricKey = crypt.ComputeSymmetricKey(crypt.UnmarshalPublicKey(pkRes.PublicKey), c.privateKey)
				log.Printf("Client %d obtained public key from server %d", c.c.Id, e.Server.Id)

				serverConn.Close()
			}

			// Send onion encrypted setup request to first hop
			data, err := crypt.EncryptOnion(b.DeadDrop, c.c.Id, c.publicKey, circuit)
			if err != nil {
				log.Fatalf("Could not onion encrypt circuit setup message at client %d", c.c.Id)
			}

			firstHopServer, serverConn := firstHop.Server.Dial()
			res, err := firstHopServer.CircuitSetup(ctx, &pb.CircuitSetupRequest{Message: data, Tag: int32(c.c.Id)})
			if err != nil {
				log.Fatalf("Circuit setup failed for client %d: %v", c.c.Id, err)
			}
			serverConn.Close()

			returnedDeadDrop, err := crypt.DecryptOnion(res.Message, circuit)
			if err != nil {
				log.Fatalf("Could not decrypt circuit setup return %d: %v", c.c.Id, err)
			}

			log.Printf("Client %d: Sent dead drop: %s, Recd dead drop: %s", c.c.Id, b.DeadDrop, returnedDeadDrop)
		}
	}
}

// TODO: split into two - send on two circuits
// func (c *client) SendMessage(ctx context.Context, msg []byte, buddy *parser.Buddy) error {
// 	encryptedMsg, err := crypt.EncryptSymmetric(msg, buddy.SymmetricKey)
// 	if err != nil {
// 		return err
// 	}
// 	ckt := buddy.Circuits[0]

// 	onionEncryptedMsg, err := crypt.EncryptOnionMessage(encryptedMsg, c.c.Id, c.publicKey, ckt.GetReversed())
// 	if err != nil {
// 		return err
// 	}
// }

func main() {
	// Setup log
	f, err := os.OpenFile("log.out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Read port from command line
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		println("Could not parse port number: %v", err)
	}

	// Initialize client
	client := client{}
	client.Initialize()

	client.c, err = parser.FetchClientByPort(port)
	if err != nil {
		log.Fatalf("Failed to get client with port: %d", port)
	}

	// Start client's server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", client.c.Host, client.c.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClientServer(s, &client)
	log.Printf("Client %d listening at %v", client.c.Id, lis.Addr())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// If client 2, add client 1 as buddy (happens symmetrically).
	// TODO: Make this fixed?
	if client.c.Id == 2 {
		log.Print("Trying to call first client")
		// Find client 1 and send symmetric key gen request
		c1 := parser.G.Clients[0]
		log.Printf("Connecting to host=%s, port=%d", c1.Host, c1.Port)

		client_1, conn := c1.Dial()
		res, err := client_1.SymmetricKeyGen(ctx, &pb.SymmetricKeyGenRequest{
			ClientId:  int32(client.c.Id),
			PublicKey: crypt.MarshalPublicKey(client.publicKey),
		})
		if err != nil {
			log.Fatalf("SymmetricKeyGen request failed: %v", err)
		}
		conn.Close()

		// Add client 1 as buddy
		client.addBuddy(1, crypt.UnmarshalPublicKey(res.PublicKey))
	}

	go s.Serve(lis)

	time.Sleep(time.Millisecond * 200) // sleep 100 ms to allow buddies to be added

	// Both clients set up their circuits
	client.setUpCircuits(ctx)

}
