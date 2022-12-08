package parser

import (
	"fmt"
	"strings"
)

type Client struct {
	Id   int    `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func FetchClientByPort(port int) (*Client, error) {
	for i, client := range G.Clients {
		if client.Port == port {
			return &G.Clients[i], nil
		}
	}

	return nil, fmt.Errorf("Client not found for port %d", port)
}

func RemoveTrailingSpaces(data []byte) []byte {
	return []byte(strings.TrimSpace(string(data)))
}
