package global

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
)

type Client struct {
	Id   int    `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Server struct {
	Id    int    `json:"id"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Layer int    `json:"layer"`
}

func (g *Global) IsTerminal(s *Server) bool {
	return s.Layer == g.FetchMaxLayers()
}

type Global struct {
	Clients []Client `json:"clients"`
	Servers []Server `json:"servers"`
}

func (g *Global) Initialize() {
	jsonFile, err := os.Open("pkg/global/global.json")
	if err != nil {
		println("Failed to parse json file: %v", err)
	}
	defer jsonFile.Close()

	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonBytes, &g)
}

func (g *Global) FetchClientByPort(port int) (*Client, error) {
	for _, client := range g.Clients {
		if client.Port == port {
			return &client, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Client not found for port %d", port))
}

func (g *Global) FetchServerByPort(port int) (*Server, error) {
	for _, server := range g.Servers {
		if server.Port == port {
			return &server, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Server not found for port %d", port))
}

func (g *Global) FetchMaxLayers() int {
	maxLayers := 0
	for _, server := range g.Servers {
		if server.Layer > maxLayers {
			maxLayers = server.Layer
		}
	}

	return maxLayers
}

func (g *Global) FetchServers(layer int) []*Server {
	var layerServers []*Server
	for _, server := range g.Servers {
		if server.Layer == layer {
			layerServers = append(layerServers, &server)
		}
	}

	return layerServers
}

func (g *Global) FetchRandomServer(layer int) *Server {
	layerServers := g.FetchServers(layer)
	i := rand.Intn(len(layerServers))

	return layerServers[i]
}

type CircuitSetupLayer struct {
	Data         []byte `json:"data"`
	EphPublicKey []byte `json:"eph_public_key"`
}

func ReadCircuitSetupLayer(data []byte) (csl *CircuitSetupLayer) {
	csl = &CircuitSetupLayer{}
	json.Unmarshal(data, csl)

	return csl
}

type CircuitForwardingData struct {
	Data           []byte `json:"data"`
	NextHopAddress string `json:"next_hop_address"`
	NextHopPort    int    `json:"next_hop_port"`
	MAC            []byte `json:"mac"`
	Tag            int    `json:"tag"`
}

func ReadCircuitForwardingData(data []byte) (cfd *CircuitForwardingData) {
	cfd = &CircuitForwardingData{}
	json.Unmarshal(data, cfd)

	return cfd
}

type CircuitElement struct {
	Server       *Server
	SymmetricKey []byte
}

type Circuit []*CircuitElement

func (c *Circuit) GetReversed() Circuit {
	ckt := Circuit{}
	for i := len(*c) - 1; i >= 0; i-- {
		ckt = append(ckt, (*c)[i])
	}

	return ckt
}

type Buddy struct {
	Id             int
	SymmetricKey   []byte
	Circuits       []Circuit
	RNG            *rand.Rand
	DeadDrop       []byte
	TerminalServer *Server
}

func (g *Global) GenerateDeadDrop(b *Buddy) {
	terminalServers := g.FetchServers(g.FetchMaxLayers())
	address := b.RNG.Int()
	i := address % len(terminalServers)
	b.DeadDrop = []byte(strconv.Itoa(address))
	b.TerminalServer = terminalServers[i]
}

func (g *Global) GenerateCircuit(b *Buddy) {
	var circuit Circuit
	maxLayers := g.FetchMaxLayers()
	for layer := 1; layer < maxLayers; layer++ {
		circuit = append(circuit, &CircuitElement{Server: g.FetchRandomServer(layer)}) // no key yet
	}
	circuit = append(circuit, &CircuitElement{Server: b.TerminalServer}) // fixed terminal server

	b.Circuits = append(b.Circuits, circuit)
}

func (g *Global) NewBuddy(id int, symmetricKey []byte) *Buddy {
	var seed int64 = 0
	for _, bt := range symmetricKey {
		seed += int64(bt)
	}
	rng := rand.New(rand.NewSource(seed))

	b := &Buddy{Id: id, SymmetricKey: symmetricKey, RNG: rng}
	g.GenerateDeadDrop(b)
	g.GenerateCircuit(b) // TODO: make two circuits

	return b
}
