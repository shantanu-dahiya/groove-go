// General structures and functions for manipulating network and node data
package parser

import (
	"encoding/json"
	"io"
	"os"
)

type Global struct {
	Clients          []Client          `json:"clients"`
	Servers          []Server          `json:"servers"`
	ServiceProviders []ServiceProvider `json:"service_providers"`
}

// These are the global variables that should be imported
// and used across the codebase to access global state
var G = Initialize()
var LastLayer = FetchMaxLayers()

// Initializes global view
func Initialize() *Global {
	jsonFile, err := os.Open("pkg/json/view.json")
	if err != nil {
		println("Failed to parse json file: %v", err)
	}
	defer jsonFile.Close()

	g := &Global{}
	jsonBytes, _ := io.ReadAll(jsonFile)
	json.Unmarshal(jsonBytes, &g)

	return g
}

// Called once to set global LastLayer
func FetchMaxLayers() int {
	maxLayers := 0
	for _, server := range G.Servers {
		if server.Layer > maxLayers {
			maxLayers = server.Layer
		}
	}

	return maxLayers
}
