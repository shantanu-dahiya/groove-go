package parser

import (
	"encoding/json"
	"fmt"
)

type MessageLayer struct {
	Data []byte `json:"data"`
	Tag  int    `json:"tag"`
}

func ReadMessageLayer(data []byte) (*MessageLayer, error) {
	ml := &MessageLayer{}
	err := json.Unmarshal(data, ml)
	if err != nil {
		fmt.Printf("Failed to parse ml: %+v\n", err)
		return nil, err
	}

	return ml, nil
}
