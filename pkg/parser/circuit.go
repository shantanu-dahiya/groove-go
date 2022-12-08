package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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

func ReadCircuitForwardingData(data []byte) (*CircuitForwardingData, error) {
	cfd := &CircuitForwardingData{}
	err := json.Unmarshal(data, cfd)
	if err != nil {
		fmt.Printf("Failed to parse cfd: %+v\n", err)
		return nil, err
	}

	return cfd, nil
}

type CircuitElement struct {
	Server       *Server
	SymmetricKey []byte
}

type Circuit []*CircuitElement

func (c *Circuit) GetReversed() *Circuit {
	ckt := &Circuit{}
	for i := len(*c) - 1; i >= 0; i-- {
		*ckt = append(*ckt, (*c)[i])
	}

	return ckt
}

func (c *Circuit) Print() string {
	out := "["
	for _, e := range *c {
		out += strconv.Itoa(e.Server.Id) + ", "
	}
	return out + "]"
}

func (b *Buddy) CreateCircuits() {
	var circuit Circuit
	for layer := 1; layer < LastLayer; layer++ {
		circuit = append(circuit, &CircuitElement{Server: FetchRandomServer(layer)})
	}
	// Terminal server dervied from dead drop ID
	circuit = append(circuit, &CircuitElement{Server: b.TerminalServer})

	b.Circuits = append(b.Circuits, &circuit)
}
