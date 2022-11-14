package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"groove-go.nyu.edu/pkg/global"
)

var ec = elliptic.P256()

type PublicKey ecdsa.PublicKey
type PrivateKey []byte

func GenerateKeys() (PublicKey, PrivateKey, error) {
	publicKey := PublicKey{}
	privateKey, x, y, err := elliptic.GenerateKey(ec, rand.Reader)
	if err != nil {
		return publicKey, privateKey, err
	}

	publicKey.X = x
	publicKey.Y = y
	publicKey.Curve = ec

	return publicKey, privateKey, nil
}

// taken from ECDH
func ComputeSymmetricKey(peerPublicKey PublicKey, privateKey PrivateKey) []byte {
	sX, _ := ec.ScalarMult(peerPublicKey.X, peerPublicKey.Y, privateKey)
	h := sha256.New()
	h.Write(sX.Bytes())
	hashedKey := h.Sum(nil)

	return hashedKey
}

func MarshalPublicKey(publicKey PublicKey) []byte {
	return elliptic.Marshal(ec, publicKey.X, publicKey.Y)
}

func UnmarshalPublicKey(data []byte) PublicKey {
	x, y := elliptic.Unmarshal(ec, data)
	return PublicKey{ec, x, y}
}

func EncryptSymmetric(message []byte, key []byte) ([]byte, error) {
	message = PKCS5Padding(message, aes.BlockSize)
	aes_cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	IV := bytes.Repeat([]byte{byte(0)}, aes_cipher.BlockSize()) // TODO: make this random

	ciphertext := make([]byte, len(message))
	aes_cipher.Encrypt(ciphertext, []byte(message))
	mode := cipher.NewCBCEncrypter(aes_cipher, IV)
	mode.CryptBlocks(ciphertext, message)

	return ciphertext, nil
}

func DecryptSymmetric(data []byte, key []byte) ([]byte, error) {
	aes_cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	IV := bytes.Repeat([]byte{byte(0)}, aes_cipher.BlockSize()) // TODO: make this random

	mode := cipher.NewCBCDecrypter(aes_cipher, IV)
	mode.CryptBlocks(data, data)

	return global.RemoveTrailingSpaces(data), err
}

// Padding is needed for AES because message length needs to be a multiple of block size.
// This is from https://gist.github.com/yingray/57fdc3264b1927ef0f984b533d63abab
func PKCS5Padding(message []byte, blockSize int) []byte {
	padding := blockSize - len(message)%blockSize
	padText := bytes.Repeat([]byte{byte('\x20')}, padding)

	return append(message, padText...)
}

func EncryptOnion(message []byte, tagId int, publicKey PublicKey, revCkt *global.Circuit) ([]byte, error) {
	// message is already encrypted with buddy's key
	for i, e := range *revCkt {
		var nextHopAddress string
		var nextHopPort int = 0

		// Since the circuit is reversed, the next hop for each server is the previous circuit element
		if i > 0 {
			nextHopAddress = (*revCkt)[i-1].Server.Host
			nextHopPort = (*revCkt)[i-1].Server.Port
		}
		cfd := &global.CircuitForwardingData{
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

		csl := &global.CircuitSetupLayer{
			Data:         encryptedCfd,
			EphPublicKey: MarshalPublicKey(publicKey),
		}

		bytes, err = json.Marshal(csl)
		if err != nil {
			return nil, err
		}

		message = bytes // set message to encrypted value for next iteration
	}

	return message, nil
}

func DecryptOnion(data []byte, ckt *global.Circuit) ([]byte, error) {
	var err error
	for _, e := range *ckt {
		data, err = DecryptSymmetric(data, e.SymmetricKey)
		if err != nil {
			return nil, err
		}

	}

	return data, nil
}

func DecryptCircuitSetupLayer(data []byte, symmetricKey []byte) (*global.CircuitForwardingData, error) {
	decryptedData, err := DecryptSymmetric(data, symmetricKey)
	if err != nil {
		return nil, err
	}

	cfd, err := global.ReadCircuitForwardingData(decryptedData)
	if err != nil {
		return nil, err
	}

	return cfd, nil
}
