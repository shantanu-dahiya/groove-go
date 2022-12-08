package parser

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Server struct {
	Id            int    `json:"id"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Layer         int    `json:"layer"`
	IsCoordinator bool   `json:"is_coordinator"`
}

func (s *Server) Print() string {
	return strconv.Itoa(s.Id)
}

func (s *Server) IsTerminal() bool {
	return s.Layer == LastLayer
}

func FetchServerByPort(port int) (*Server, error) {
	for i, server := range G.Servers {
		if server.Port == port {
			return &G.Servers[i], nil
		}
	}

	return nil, fmt.Errorf("Server not found for port %d", port)
}

func FetchServers(layer int) []*Server {
	var layerServers []*Server
	for i, s := range G.Servers {
		if s.Layer == layer {
			// Fetch by index - reference doesn't work
			layerServers = append(layerServers, &G.Servers[i])
		}
	}

	return layerServers
}

func FetchRandomServer(layer int) *Server {
	layerServers := FetchServers(layer)
	i := rand.Intn(len(layerServers))

	return layerServers[i]
}
