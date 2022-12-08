package crypt

import (
	"bytes"
	"encoding/json"

	parser "groove-go.nyu.edu/pkg/parser"
)

// Padding is needed for AES because message length needs to be a multiple of block size.
// This is from https://gist.github.com/yingray/57fdc3264b1927ef0f984b533d63abab
func PKCS5Padding(message []byte, blockSize int) []byte {
	padding := blockSize - len(message)%blockSize
	padText := bytes.Repeat([]byte{byte('\x20')}, padding)

	return append(message, padText...)
}

// Onion encrypt circuit setup message
func EncryptOnion(message []byte, tagId int, publicKey PublicKey, ckt *parser.Circuit) ([]byte, error) {
	revCkt := ckt.GetReversed()
	// Message is dead drop encrypted with buddy's key
	for i, e := range *revCkt {
		var nextHopAddress string
		var nextHopPort int = 0

		// Since the circuit is reversed, the next hop for each server is the previous circuit element
		if i > 0 {
			nextHopAddress = (*revCkt)[i-1].Server.Host
			nextHopPort = (*revCkt)[i-1].Server.Port
		}
		cfd := &parser.CircuitForwardingData{
			Data:           message,
			NextHopAddress: nextHopAddress,
			NextHopPort:    nextHopPort,
			MAC:            nil,
			Tag:            tagId,
		}

		bytes, err := json.Marshal(cfd)
		if err != nil {
			return nil, err
		}

		encryptedCfd, err := EncryptSymmetric(bytes, e.SymmetricKey)
		if err != nil {
			return nil, err
		}

		csl := &parser.CircuitSetupLayer{
			Data:         encryptedCfd,
			EphPublicKey: MarshalPublicKey(publicKey),
		}

		bytes, err = json.Marshal(csl)
		if err != nil {
			return nil, err
		}

		// Set message to encrypted value for next iteration
		message = bytes
	}

	return message, nil
}

func EncryptOnionMessage(message []byte, tagId int, publicKey PublicKey, ckt *parser.Circuit) ([]byte, error) {
	revCkt := ckt.GetReversed()
	for _, e := range *revCkt {
		ml := &parser.MessageLayer{
			Data: message,
			Tag:  tagId,
		}

		bytes, err := json.Marshal(ml)
		if err != nil {
			return nil, err
		}

		bytes, err = EncryptSymmetric(bytes, e.SymmetricKey)
		if err != nil {
			return nil, err
		}

		// Set message to encrypted value for next iteration
		message = bytes
	}

	return message, nil
}

func DecryptOnion(data []byte, ckt *parser.Circuit) ([]byte, error) {
	var err error
	for _, e := range *ckt {
		data, err = DecryptSymmetric(data, e.SymmetricKey)
		if err != nil {
			return nil, err
		}

	}

	return data, nil
}

func DecryptCircuitSetupLayer(data []byte, symmetricKey []byte) (*parser.CircuitForwardingData, error) {
	decryptedData, err := DecryptSymmetric(data, symmetricKey)
	if err != nil {
		return nil, err
	}

	cfd, err := parser.ReadCircuitForwardingData(decryptedData)
	if err != nil {
		return nil, err
	}

	return cfd, nil
}
